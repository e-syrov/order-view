package nats

import (
	"fmt"
	"github.com/nats-io/stan.go"
	"log"
	"wb/internal/cache"
	"wb/internal/database/psql"
	"wb/internal/models"
)

const OrdersChan = "orders-channel"

type NatsInterface interface {
	SubscribeOnChannel(nameCh string, csm *cache.CacheSM, db *psql.Database, errorChannel chan error)
	PublishOnChannel(nameCh string, message []byte)
}

type Nats struct {
	Conn stan.Conn
}

func New() *Nats {
	return &Nats{}
}
func (nats *Nats) Connect(clusterID, clientID string, errorChannel chan error) {
	var err error
	nats.Conn, err = stan.Connect(clusterID, clientID)
	if err != nil {
		errorChannel <- fmt.Errorf("ошибка подключения к NATS Streaming: %v", err)
	}
}

func (nats Nats) SubscribeOnChannel(nameCh string, csm *cache.CacheSM, db *psql.Database, errorChannel chan error) {
	_, err := nats.Conn.Subscribe(nameCh, func(msg *stan.Msg) {
		var ord models.Order
		err := ord.FromJSON(msg.Data)
		if err != nil {
			log.Printf("ошибка декодирования JSON: %v", err)
			return
		}
		if ord.OrderUid == "" {
			log.Printf("получен некорректный order_uid (%s)", ord.OrderUid)
			return
		}

		if _, ok := csm.Cache.Load(ord.OrderUid); !ok {
			err = csm.SaveToCache(ord)
			if err != nil {
				log.Printf("ошибка сохранения в кэш: %v", err)
			}
			err = db.SaveToDB(ord)
			if err != nil {
				errorChannel <- fmt.Errorf("ошибка сохранения в БД: %v", err)
			}
		} else {
			log.Printf("order с order_uid = %s уже существует\n", ord.OrderUid)
		}

	})
	if err != nil {
		errorChannel <- fmt.Errorf("ошибка подписки: %v ", err)
	}
}

func (nats Nats) PublishOnChannel(nameCh string, message []byte) {
	err := nats.Conn.Publish(nameCh, message)
	if err != nil {
		log.Fatalf("Ошибка публикации: %v", err)
	}
}
