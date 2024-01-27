package database

import (
	"database/sql"
	"fmt"

	model "example.com/service/service/internal/models"
	"github.com/spf13/viper"
)

type Database struct {
	db *sql.DB
}

func Initialize() (*Database, error) {
	viper.AddConfigPath("internal\\configs")
	viper.SetConfigName("config")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	connect_db := "host=" + viper.GetString("db.host") + " " + "user=" + viper.GetString("db.username") + " " + "port=" + viper.GetString("db.port") + " " + "password=" + viper.GetString("db.password") + " " + "dbname=" + viper.GetString("db.dbname") + " " + "sslmode=" + viper.GetString("db.sslmode")
	db, err := sql.Open("postgres", connect_db)
	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &Database{db: db}, nil
}

func (db *Database) GetAllOrdersDatabase() ([]model.OrderDetails, error) {
	query := `
        SELECT 
            o.order_uid, o.track_number, o.entry, o.locale, o.internal_signature,
            o.customer_id, o.delivery_service, o.shardkey, o.sm_id, o.date_created,
            o.oof_shard,
            d.name as delivery_name, d.phone as delivery_phone, d.zip, d.city, d.address,
            d.region, d.email,
            p.transaction, p.request_id, p.currency, p.provider, p.amount, p.payment_dt,
            p.bank, p.delivery_cost, p.goods_total, p.custom_fee
        FROM orders o
        LEFT JOIN delivery d ON o.order_uid = d.order_uid
        LEFT JOIN payment p ON o.order_uid = p.order_uid
    `

	rows, err := db.db.Query(query)
	if err != nil {
		fmt.Println("Error executing query:", err)
		return nil, err
	}
	defer rows.Close()

	var ordDet []model.OrderDetails
	var orderDetails model.OrderDetails

	for rows.Next() {
		err := rows.Scan(
			&orderDetails.OrderUID, &orderDetails.TrackNumber, &orderDetails.Entry, &orderDetails.Locale,
			&orderDetails.InternalSignature, &orderDetails.CustomerID, &orderDetails.DeliveryService,
			&orderDetails.ShardKey, &orderDetails.SMID, &orderDetails.DateCreated, &orderDetails.OOFShard,
			&orderDetails.Delivery.Name, &orderDetails.Delivery.Phone, &orderDetails.Delivery.Zip,
			&orderDetails.Delivery.City, &orderDetails.Delivery.Address, &orderDetails.Delivery.Region,
			&orderDetails.Delivery.Email, &orderDetails.Payment.Transaction, &orderDetails.Payment.RequestID,
			&orderDetails.Payment.Currency, &orderDetails.Payment.Provider, &orderDetails.Payment.Amount,
			&orderDetails.Payment.PaymentDt, &orderDetails.Payment.Bank, &orderDetails.Payment.DeliveryCost,
			&orderDetails.Payment.GoodsTotal, &orderDetails.Payment.CustomFee,
		)

		if err != nil {
			fmt.Println("Error scanning row:", err)
			return []model.OrderDetails{}, err
		}

		ordDet = append(ordDet, orderDetails)
	}

	if err := rows.Err(); err != nil {
		fmt.Println("Error iterating over rows:", err)
		return []model.OrderDetails{}, err
	}

	for i := 0; i < len(ordDet); i++ {
		query := `
        	SELECT 
				chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status
			FROM items
			WHERE order_uid = $1
    	`

		rows, err := db.db.Query(query, ordDet[i].OrderUID)
		if err != nil {
			fmt.Println("Error executing query:", err)
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var item model.Item
			err := rows.Scan(
				&item.ChrtID, &item.TrackNumber, &item.Price, &item.RID, &item.Name, &item.Sale,
				&item.Size, &item.TotalPrice, &item.NMID, &item.Brand, &item.Status,
			)

			if err != nil {
				fmt.Println("Error scanning row:", err)
				return []model.OrderDetails{}, err
			}

			ordDet[i].Items = append(ordDet[i].Items, item)
		}
	}

	return ordDet, nil
}

func (db *Database) GetDatabase(id string) (model.OrderDetails, error) {
	query := `
		SELECT 
			o.order_uid, o.track_number, o.entry, o.locale, o.internal_signature,
			o.customer_id, o.delivery_service, o.shardkey, o.sm_id, o.date_created,
			o.oof_shard,
			d.name as delivery_name, d.phone as delivery_phone, d.zip, d.city, d.address,
			d.region, d.email,
			p.transaction, p.request_id, p.currency, p.provider, p.amount, p.payment_dt,
			p.bank, p.delivery_cost, p.goods_total, p.custom_fee,
			i.chrt_id, i.track_number, i.price, i.rid, i.name as item_name, i.sale, i.size,
			i.total_price, i.nm_id, i.brand, i.status
		FROM orders o
		LEFT JOIN delivery d ON o.order_uid = d.order_uid
		LEFT JOIN payment p ON o.order_uid = p.order_uid
		LEFT JOIN items i ON o.order_uid = i.order_uid
		WHERE o.order_uid = $1
	`

	rows, err := db.db.Query(query, id)
	if err != nil {
		fmt.Println("Error executing query:", err)
		return model.OrderDetails{}, err
	}
	defer rows.Close()

	var orderDetails model.OrderDetails
	orderDetails.Items = []model.Item{}

	for rows.Next() {
		var item model.Item
		err := rows.Scan(
			&orderDetails.OrderUID, &orderDetails.TrackNumber, &orderDetails.Entry, &orderDetails.Locale,
			&orderDetails.InternalSignature, &orderDetails.CustomerID, &orderDetails.DeliveryService,
			&orderDetails.ShardKey, &orderDetails.SMID, &orderDetails.DateCreated, &orderDetails.OOFShard,
			&orderDetails.Delivery.Name, &orderDetails.Delivery.Phone, &orderDetails.Delivery.Zip,
			&orderDetails.Delivery.City, &orderDetails.Delivery.Address, &orderDetails.Delivery.Region,
			&orderDetails.Delivery.Email, &orderDetails.Payment.Transaction, &orderDetails.Payment.RequestID,
			&orderDetails.Payment.Currency, &orderDetails.Payment.Provider, &orderDetails.Payment.Amount,
			&orderDetails.Payment.PaymentDt, &orderDetails.Payment.Bank, &orderDetails.Payment.DeliveryCost,
			&orderDetails.Payment.GoodsTotal, &orderDetails.Payment.CustomFee, &item.ChrtID, &item.TrackNumber, &item.Price,
			&item.RID, &item.Name, &item.Sale, &item.Size, &item.TotalPrice, &item.NMID, &item.Brand,
			&item.Status,
		)

		if err != nil {
			fmt.Println("Error scanning row:", err)
			return model.OrderDetails{}, err
		}

		orderDetails.Items = append(orderDetails.Items, item)
	}

	if err := rows.Err(); err != nil {
		fmt.Println("Error iterating over rows:", err)
		return model.OrderDetails{}, err
	}

	return orderDetails, nil
}

func (db *Database) SendingDatabase(orderDetails model.OrderDetails) error {

	// Start a transaction
	tx, err := db.db.Begin()
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
			INSERT INTO Items (order_uid, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
			orderDetails.Order.OrderUID, item.ChrtID, item.TrackNumber, item.Price, item.RID, item.Name, item.Sale, item.Size, item.TotalPrice, item.NMID, item.Brand, item.Status)
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

func (db *Database) Close() error {
	return db.db.Close()
}
