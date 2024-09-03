CREATE DATABASE balance-service;

--  пользователи и их баланс
create table users
(
    id      INT             PRIMARY KEY,
    balance NUMERIC(10, 2)  NOT NULL DEFAULT 0 CHECK (balance >= 0)
);


-- резервирование средств с основого баланса на отдельном счете
create table transactions
(
    id         SERIAL PRIMARY KEY,
    user_id    INTEGER        NOT NULL,
    service_id INTEGER        NOT NULL DEFAULT 0,
    order_id   INTEGER        NOT NULL DEFAULT 0,
    amount     NUMERIC(10, 2) NOT NULL DEFAULT 0 CHECK (amount >= 0),
    transaction_type     VARCHAR(10)    NOT NULL, -- deposit, reserve, reserve_confirm, reserve_cancel

        FOREIGN KEY (user_id) REFERENCES users(id)
);


-- отчет для бухгалтерии
create table report
(
    id         SERIAL PRIMARY KEY,
    service_id INTEGER        NOT NULL,
    amount     NUMERIC(10, 2) NOT NULL DEFAULT 0 CHECK (amount >= 0)
);
