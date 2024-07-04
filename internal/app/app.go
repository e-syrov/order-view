package app

import (
	"fmt"
	"net/http"
	"sync"
	"wb/internal/cache"
	"wb/internal/database/psql"
	"wb/internal/nats"
	handlers "wb/internal/web"
)

var (
	wg           = &sync.WaitGroup{}
	cacheOrders  = cache.New()
	webServer    = handlers.New(&cacheOrders)
	db           = psql.New()
	ErrorChannel = make(chan error)
	ns           = nats.New()
)

func Run() error {

	err := db.InitDB()
	if err != nil {
		return err
	}
	fmt.Println("[DB]Подключение к бд завершено")

	orders, err := db.GetAllOrders()
	if err != nil {
		return err
	}
	if orders == nil {
		fmt.Println("[DB]В бд нет записей orders")
	} else {
		fmt.Println("[DB]Получение всех orders из бд завершено")

		for _, order := range orders {
			err = cacheOrders.SaveToCache(order)
			if err != nil {
				return err
			}
		}
		fmt.Println("[Cache]Восстановление кэша из бд завершено")
	}

	wg.Add(1)
	go func() {
		defer wg.Done()

		ns.Connect("test-cluster", "subscriber", ErrorChannel)
		fmt.Println("[NATS]Подключение к NATS Streaming завершено")

		ns.SubscribeOnChannel(nats.OrdersChan, &cacheOrders, db, ErrorChannel)
		fmt.Printf("[NATS]Подписка на канал %s активна\n", nats.OrdersChan)

	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		http.HandleFunc("/order", webServer.HandlerOrders)
		err = http.ListenAndServe(":8080", nil)
		if err != nil {
			ErrorChannel <- fmt.Errorf("ошибка начала работы сервера: %v", err)
		}
		fmt.Println("[WEB]Начало работы сервера на порту 8080")
	}()

	select {
	case err := <-ErrorChannel:
		return err
	case <-waitForWaitGroup(wg):
		return nil
	}

}

func waitForWaitGroup(wg *sync.WaitGroup) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()
	return done
}
