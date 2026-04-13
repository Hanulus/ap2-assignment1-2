-- Orders table for the Order Service
-- Run this script once before starting the service

CREATE TABLE IF NOT EXISTS orders (
    id          VARCHAR(36) PRIMARY KEY,
    customer_id VARCHAR(100) NOT NULL,
    item_name   VARCHAR(255) NOT NULL,
    amount      BIGINT NOT NULL,        -- stored in cents
    status      VARCHAR(20) NOT NULL,   -- Pending, Paid, Failed, Cancelled
    created_at  TIMESTAMP NOT NULL
);
