-- Create the 'orders' table
CREATE TABLE orders
(
    id         BIGINT PRIMARY KEY AUTO_INCREMENT,
    status     VARCHAR(50)    NOT NULL,
    amount     DECIMAL(10, 2) NOT NULL,
    created_at DATETIME       NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME       NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Create the 'payments' table
CREATE TABLE payments
(
    id           BIGINT PRIMARY KEY AUTO_INCREMENT,
    status       VARCHAR(50)    NOT NULL,
    order_id     BIGINT         NOT NULL,
    payment_type VARCHAR(50)    NOT NULL,
    amount       DECIMAL(10, 2) NOT NULL,
    details      VARCHAR(200),
    created_at   DATETIME       NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   DATETIME       NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT fk_payments_order
        FOREIGN KEY (order_id) REFERENCES orders (id)
            ON DELETE CASCADE
);

-- Create the 'charges' table
CREATE TABLE charges
(
    id         BIGINT PRIMARY KEY AUTO_INCREMENT,
    amount     DECIMAL(10, 2) NOT NULL,
    category   VARCHAR(50)    NOT NULL,
    payment_id BIGINT         NOT NULL,
    created_at DATETIME       NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME       NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT fk_charges_payment
        FOREIGN KEY (payment_id) REFERENCES payments (id)
            ON DELETE CASCADE
);

-- Insert sample data into 'orders' table
INSERT INTO orders (status, amount)
VALUES ('pending', 120.50),
       ('pending', 250.00),
       ('pending', 89.99),
       ('pending', 310.25),
       ('pending', 45.00),
       ('pending', 199.99),
       ('pending', 540.75),
       ('pending', 123.45),
       ('pending', 79.90),
       ('pending', 65.25),
       ('pending', 300.00),
       ('pending', 410.10),
       ('pending', 145.50),
       ('pending', 275.80),
       ('pending', 88.88),
       ('pending', 59.99),
       ('pending', 160.60),
       ('pending', 330.33),
       ('pending', 200.00),
       ('pending', 110.10),
       ('pending', 500.00),
       ('pending', 215.25),
       ('pending', 39.99),
       ('pending', 90.00),
       ('pending', 149.49);
