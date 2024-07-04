package psql

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"wb/internal/models"
)

var (
	db *sql.DB
)

const (
	host     = "localhost"
	port     = 5432
	user     = "hello"
	password = "hello"
	dbname   = "ordersdb"
)

type DBInterface interface {
	InitDB() error
	GetAllOrders() ([]models.Order, error)
	SaveToDB(order models.Order) error
}

type Database struct {
	DB *sql.DB
}

func New() *Database {
	return &Database{}
}

func (db *Database) InitDB() error {
	connStr := fmt.Sprintf("user=%s dbname=%s sslmode=disable password=%s", user, dbname, password)
	var err error
	db.DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("ошибка подключения к базе данных: %v", err)
	}
	return nil
}

func (db *Database) GetAllOrders() ([]models.Order, error) {
	var err error
	queryOrders := `SELECT 
    o.order_uid, o.track_number, o.entry, o.locale, o.internal_signature, 
    o.customer_id, o.delivery_service, o.shardkey, o.sm_id, o.date_created, o.oof_shard,
    d.name, d.phone, d.zip, d.city, d.address, d.region, d.email,
    p.transaction, p.request_id, p.currency, p.provider, p.amount, 
    p.payment_dt, p.bank, p.delivery_cost, p.goods_total, p.custom_fee
FROM 
    orders o
LEFT JOIN 
    deliveries d ON o.order_uid = d.order_uid
LEFT JOIN 
    payments p ON o.order_uid = p.order_uid;`

	ordersRows, err := db.DB.Query(queryOrders)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения данных(ordersRows): %v", err)
	}
	defer ordersRows.Close()

	var orders []models.Order

	for ordersRows.Next() {
		var ord models.Order
		err = ordersRows.Scan(
			&ord.OrderUid, &ord.TrackNumber, &ord.Entry, &ord.Locale, &ord.InternalSignature,
			&ord.CustomerID, &ord.DeliveryService, &ord.Shardkey, &ord.SmID, &ord.DateCreated, &ord.OofShard,
			&ord.Delivery.Name, &ord.Delivery.Phone, &ord.Delivery.Zip, &ord.Delivery.City, &ord.Delivery.Address, &ord.Delivery.Region, &ord.Delivery.Email,
			&ord.Payment.Transaction, &ord.Payment.RequestID, &ord.Payment.Currency, &ord.Payment.Provider, &ord.Payment.Amount,
			&ord.Payment.PaymentDt, &ord.Payment.Bank, &ord.Payment.DeliveryCost, &ord.Payment.GoodsTotal, &ord.Payment.CustomFee)
		if err != nil {
			return nil, fmt.Errorf("ошибка чтения данных(ordersRows): %v", err)
		}
		queryItems := `SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status FROM items WHERE order_uid = $1`

		itemsRows, err := db.DB.Query(queryItems, ord.OrderUid)
		if err != nil {
			return nil, fmt.Errorf("ошибка получения данных(itemsRows): %v", err)
		}

		for itemsRows.Next() {
			var it models.Item
			err := itemsRows.Scan(&it.ChrtID, &it.TrackNumber, &it.Price, &it.Rid, &it.Name, &it.Sale, &it.Size, &it.TotalPrice, &it.NmID, &it.Brand, &it.Status)
			if err != nil {
				return nil, fmt.Errorf("ошибка чтения данных(itemsRows): %v", err)
			}
			ord.Items = append(ord.Items, it)
		}
		orders = append(orders, ord)

	}
	return orders, nil

}

func (db *Database) SaveToDB(order models.Order) error {
	tx, err := db.DB.Begin()
	if err != nil {
		return fmt.Errorf("ошибка начала транзакции: %v", err)
	}
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Восстановление после паники: %v", r)
			if err := tx.Rollback(); err != nil {
				log.Fatalf("Ошибка отката транзакции после паники: %v", err)
			}
		}
	}()

	queryDeliveries := `INSERT INTO deliveries (order_uid, name, phone, zip, city, address, region, email)
							VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	queryPayments := `INSERT INTO payments (order_uid, transaction, request_id, currency, provider, amount, 
                    	  payment_dt, bank, delivery_cost, goods_total, custom_fee ) 
						  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	queryItems := `INSERT INTO items (order_uid, chrt_id, track_number, price, rid, name, sale, size, 
                   total_price, nm_id, brand, status)
                   VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	queryOrders := `INSERT INTO orders (order_uid, track_number, entry, locale, internal_signature, 
                    customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard) 
                    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	_, err = tx.Exec(queryOrders, order.OrderUid, order.TrackNumber, order.Entry, order.Locale,
		order.InternalSignature, order.CustomerID, order.DeliveryService, order.Shardkey,
		order.SmID, order.DateCreated, order.OofShard)
	if err != nil {
		return fmt.Errorf("ошибка добавления данных в таблицу Orders: %v", err)
	}

	_, err = tx.Exec(queryDeliveries, order.OrderUid, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip,
		order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email)

	if err != nil {
		return fmt.Errorf("ошибка добавления данных в таблицу Deliveries: %v", err)
	}

	_, err = tx.Exec(queryPayments, order.OrderUid, order.Payment.Transaction, order.Payment.RequestID, order.Payment.Currency, order.Payment.Provider,
		order.Payment.Amount, order.Payment.PaymentDt, order.Payment.Bank, order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee)
	if err != nil {
		return fmt.Errorf("ошибка добавления данных в таблицу Payments: %v", err)
	}

	for _, it := range order.Items {
		_, err = tx.Exec(queryItems, order.OrderUid, it.ChrtID, it.TrackNumber, it.Price, it.Rid, it.Name,
			it.Sale, it.Size, it.TotalPrice, it.NmID, it.Brand, it.Status)
		if err != nil {
			return fmt.Errorf("ошибка добавления данных в таблицу Items: %v", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("ошибка подтверждения изменений: %v", err)
	}
	return err
}
