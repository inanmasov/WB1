package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"time"

	model "example.com/service/service/internal/models"
	_ "github.com/lib/pq"
	"github.com/nats-io/stan.go"
	"github.com/spf13/viper"
)

func main() {

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

		var order model.OrderDetails

		err := json.Unmarshal([]byte(msg.Data), &order)
		if err != nil {
			log.Fatal(err)
		}
		if err := insertOrderDetails(order); err != nil {
			log.Fatal(err)
		}
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

func insertOrderDetails(orderDetails model.OrderDetails) error {
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

	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Insert into Orders table
	_, err = tx.Exec(`
		INSERT INTO Orders (order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		orderDetails.Order.OrderUID, orderDetails.Order.TrackNumber, orderDetails.Order.Entry, orderDetails.Order.Locale,
		orderDetails.Order.InternalSignature, orderDetails.Order.CustomerID, orderDetails.Order.DeliveryService, orderDetails.Order.ShardKey,
		orderDetails.Order.SMID, orderDetails.Order.DateCreated, orderDetails.Order.OOFShard)
	if err != nil {
		return err
	}

	// Insert into Delivery table
	_, err = tx.Exec(`
		INSERT INTO Delivery (order_uid, name, phone, zip, city, address, region, email)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		orderDetails.Order.OrderUID, orderDetails.Delivery.Name, orderDetails.Delivery.Phone, orderDetails.Delivery.Zip,
		orderDetails.Delivery.City, orderDetails.Delivery.Address, orderDetails.Delivery.Region, orderDetails.Delivery.Email)
	if err != nil {
		return err
	}

	// Insert into Payment table
	_, err = tx.Exec(`
		INSERT INTO Payment (order_uid, transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		orderDetails.Order.OrderUID, orderDetails.Payment.Transaction, orderDetails.Payment.RequestID, orderDetails.Payment.Currency,
		orderDetails.Payment.Provider, orderDetails.Payment.Amount, orderDetails.Payment.PaymentDt, orderDetails.Payment.Bank,
		orderDetails.Payment.DeliveryCost, orderDetails.Payment.GoodsTotal, orderDetails.Payment.CustomFee)
	if err != nil {
		return err
	}

	// Insert into Items table
	for _, item := range orderDetails.Items {
		_, err := tx.Exec(`
			INSERT INTO Items (order_uid, chrt_id, price, rid, name, sale, size, total_price, nm_id, brand, status)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
			orderDetails.Order.OrderUID, item.ChrtID, item.Price, item.RID, item.Name, item.Sale, item.Size, item.TotalPrice, item.NMID, item.Brand, item.Status)
		if err != nil {
			return err
		}
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
