package main

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
	"github.com/nats-io/stan.go"
	"github.com/spf13/viper"
)

func main() {
	viper.AddConfigPath("service\\configs")
	viper.SetConfigName("config")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	connect_db := "host=" + viper.GetString("db.host") + " " + "user=" + viper.GetString("db.username") + " " + "port=" + viper.GetString("db.port") + " " + "password=" + viper.GetString("db.password") + " " + "dbname=" + viper.GetString("db.dbname") + " " + "sslmode=" + viper.GetString("db.sslmode")
	db, err := sql.Open("postgres", connect_db)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Подключение к серверу NATS Streaming
	sc, err := stan.Connect("test-cluster", "publisher-client", stan.NatsURL("nats://localhost:4222"))
	if err != nil {
		log.Fatalf("Ошибка подключения к серверу NATS Streaming: %v", err)
	}
	defer sc.Close()

	// Название канала, на который вы хотите подписаться
	channel := "myrad"

	// Функция, которая будет вызвана при получении сообщения
	handler := func(msg *stan.Msg) {
		log.Printf("Получено сообщение из канала %s: %s", channel, string(msg.Data))
	}

	// Подписка на канал
	subscription, err := sc.Subscribe(channel, handler, stan.DeliverAllAvailable())
	if err != nil {
		log.Fatalf("Ошибка подписки на канал: %v", err)
	}
	defer subscription.Unsubscribe()

	log.Printf("Подписан на канал %s", channel)
	time.Sleep(3 * time.Second)
}
