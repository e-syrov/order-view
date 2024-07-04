package cache

import (
	"fmt"
	"sync"
	"wb/internal/models"
)

// CacheInterface определяет методы для сохранения и извлечения заказов из кэша
type CacheInterface interface {
	SaveToCache(order models.Order) error
	GetFromCache(orderUid string) (models.Order, error)
}

// CacheSM реализует интерфейс CacheInterface с использованием sync.Map для потокобезопасного хранения данных
type CacheSM struct {
	Cache sync.Map
}

// NewCacheSM создает новый экземпляр CacheSM
func New() CacheSM {
	return CacheSM{
		Cache: sync.Map{},
	}
}

// SaveToCache сохраняет заказ в кэш, если заказ с таким orderUid еще не существует
func (csm *CacheSM) SaveToCache(order models.Order) error {
	if _, found := csm.Cache.Load(order.OrderUid); !found {
		csm.Cache.Store(order.OrderUid, order)
		return nil
	}
	return fmt.Errorf("order с order_uid = %v уже существует", order.OrderUid)
}

// GetFromCache извлекает заказ из кэша по orderUid
func (csm *CacheSM) GetFromCache(orderUid string) (models.Order, error) {
	ord, found := csm.Cache.Load(orderUid)
	if found {
		return ord.(models.Order), nil
	}
	return models.Order{}, fmt.Errorf("order с order_uid = %v не найден", orderUid)
}
