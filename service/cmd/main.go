package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"example.com/service/service/internal/cached"
	database "example.com/service/service/internal/database"
	model "example.com/service/service/internal/models"
	reg "example.com/service/service/internal/transport/router"
	_ "github.com/lib/pq"
	"github.com/nats-io/stan.go"
	"github.com/spf13/viper"
)

func main() {
	// Cache init
	//cacheManager := cached.NewCache(1*time.Hour, 24*time.Hour)
	// Подключение к серверу NATS Streaming
	sc, err := stan.Connect("test-cluster", "publisher-client", stan.NatsURL("nats://localhost:4222"))
	if err != nil {
		log.Fatalf("Ошибка подключения к серверу NATS Streaming: %v", err)
	}
	defer sc.Close()

	// Название канала, на который вы хотите подписаться
	channel := "myrad"

	var wg sync.WaitGroup

	flag := cached.InitCacheDB()

	// Функция, которая будет вызвана при получении сообщения
	handler := func(msg *stan.Msg) {
		defer wg.Done()

		var order model.OrderDetails

		err := json.Unmarshal([]byte(msg.Data), &order)
		if err != nil {
			log.Fatal(err)
		}
		db, err := database.Initialize()
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		if !flag {
			key := order.OrderUID
			cached.GlobalCacheManager.Cache.Set(key, order, -1)

			if err := db.SendingDatabase(order); err != nil {
				log.Fatal(err)
			}
		}
	}

	// Подписка на канал
	subscription, err := sc.Subscribe(channel, func(msg *stan.Msg) {
		wg.Add(1)
		go handler(msg)
	}, stan.DeliverAllAvailable())

	if err != nil {
		log.Fatalf("Ошибка подписки на канал: %v", err)
	}
	defer subscription.Unsubscribe()
	log.Printf("Подписан на канал %s", channel)
	time.Sleep(time.Second)

	reg.RegisterRoutes()

	viper.AddConfigPath("service\\internal\\configs")
	viper.SetConfigName("config")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	log.Println("server start listening on port", viper.GetString("port"))
	err = http.ListenAndServe(":"+viper.GetString("port"), nil)
	if err != nil {
		log.Fatal(err)
	}
}
