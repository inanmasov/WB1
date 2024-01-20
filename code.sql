DROP TABLE IF EXISTS orders CASCADE;
DROP TABLE IF EXISTS delivery CASCADE;
DROP TABLE IF EXISTS payment CASCADE;
DROP TABLE IF EXISTS items CASCADE;


CREATE TABLE IF NOT EXISTS orders (
    order_uid VARCHAR(255) PRIMARY KEY,
    track_number VARCHAR(255),
    entry VARCHAR(255),
    locale VARCHAR(5),
    internal_signature VARCHAR(255),
    customer_id VARCHAR(255),
    delivery_service VARCHAR(255),
    shardkey VARCHAR(10),
    sm_id INT,
    date_created TIMESTAMP,
    oof_shard VARCHAR(10)
);

CREATE TABLE IF NOT EXISTS delivery (
    order_uid VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255),
    phone VARCHAR(15),
    zip VARCHAR(20),
    city VARCHAR(255),
    address VARCHAR(255),
    region VARCHAR(255),
    email VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS payment (
    order_uid VARCHAR(255) PRIMARY KEY,
    transaction VARCHAR(255),
    request_id VARCHAR(255),
    currency VARCHAR(5),
    provider VARCHAR(255),
    amount INT,
    payment_dt INT,
    bank VARCHAR(255),
    delivery_cost INT,
    goods_total INT,
    custom_fee INT
);

CREATE TABLE IF NOT EXISTS items (
    order_uid VARCHAR(255),
    chrt_id INT,
    price INT,
    rid VARCHAR(255),
    name VARCHAR(255),
    sale INT,
    size VARCHAR(10),
    total_price INT,
    nm_id INT,
    brand VARCHAR(255),
    status INT,
    PRIMARY KEY (order_uid, chrt_id),
    FOREIGN KEY (order_uid) REFERENCES orders (order_uid)
);



-- Insert into orders table
INSERT INTO orders (order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
VALUES ('b563feb7b2b84b6test123', 'WBILMTESTTRACKqwe', 'WBIL1', 'en', '', 'test', 'meest', '9', 99, '2021-11-26T06:22:19Z', '1');

-- Insert into delivery_info table
INSERT INTO delivery (order_uid, name, phone, zip, city, address, region, email)
VALUES ('b563feb7b2b84b6test123', 'Test Testov qwe', '+19720000000', '26398092', 'Kiryat Mozkin3', 'Ploshad Mira 152', 'Kraiote', 'test1@gmail.com');

-- Insert into payment_info table
INSERT INTO payment (order_uid, transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
VALUES ('b563feb7b2b84b6test123', 'b563feb7b2b84b6test123', '', 'USD', 'wbpay', 18171, 16379077272, 'alpha123', 15001, 3173, 0);

-- Insert into order_items table (first item)
INSERT INTO items (order_uid, chrt_id, price, rid, name, sale, size, total_price, nm_id, brand, status)
VALUES ('b563feb7b2b84b6test123', 1234, 453, 'ab4219087a764ae0btest123', 'Mascaras1', 30, '0', 317, 2389212, 'Vivienne Sabo123', 202);

-- Insert into order_items table (second item)
INSERT INTO items (order_uid, chrt_id, price, rid, name, sale, size, total_price, nm_id, brand, status)
VALUES ('b563feb7b2b84b6test123', 12345, 4531, 'ab4219087a764ae0btest123a', 'Mascaras1s', 30, '01', 3171, 23892122, 'Vivienne Sabo1233', 202);
