-- Основная таблица заказов
CREATE TABLE orders (
                        order_uid TEXT PRIMARY KEY,
                        track_number TEXT NOT NULL,
                        entry TEXT NOT NULL,
                        locale TEXT,
                        internal_signature TEXT,
                        customer_id TEXT NOT NULL,
                        delivery_service TEXT,
                        shardkey TEXT,
                        sm_id INT,
                        date_created TIMESTAMP NOT NULL,
                        oof_shard TEXT
);

-- Доставка
CREATE TABLE delivery (
                          id SERIAL PRIMARY KEY,
                          order_uid TEXT NOT NULL REFERENCES orders(order_uid) ON DELETE CASCADE,
                          name TEXT NOT NULL,
                          phone TEXT,
                          zip TEXT,
                          city TEXT NOT NULL,
                          address TEXT NOT NULL,
                          region TEXT,
                          email TEXT
);

-- Оплата
CREATE TABLE payment (
                         id SERIAL PRIMARY KEY,
                         order_uid TEXT NOT NULL REFERENCES orders(order_uid) ON DELETE CASCADE,
                         transaction TEXT NOT NULL,
                         request_id TEXT,
                         currency TEXT NOT NULL,
                         provider TEXT,
                         amount INT NOT NULL,
                         payment_dt BIGINT NOT NULL,
                         bank TEXT,
                         delivery_cost INT,
                         goods_total INT,
                         custom_fee INT
);

-- Товары
CREATE TABLE items (
                       id SERIAL PRIMARY KEY,
                       order_uid TEXT NOT NULL REFERENCES orders(order_uid) ON DELETE CASCADE,
                       chrt_id INT NOT NULL,
                       track_number TEXT NOT NULL,
                       price INT NOT NULL,
                       rid TEXT NOT NULL,
                       name TEXT NOT NULL,
                       sale INT,
                       size TEXT,
                       total_price INT NOT NULL,
                       nm_id INT,
                       brand TEXT,
                       status INT
);

-- Индексы для ускорения поиска
CREATE INDEX idx_orders_date_created ON orders(date_created);
CREATE INDEX idx_payment_transaction ON payment(transaction);
CREATE INDEX idx_items_chrt_id ON items(chrt_id);
