-- Payments table for the Payment Service
-- Run this script once before starting the service

CREATE TABLE IF NOT EXISTS payments (
    id             VARCHAR(36) PRIMARY KEY,
    order_id       VARCHAR(36) NOT NULL UNIQUE,
    transaction_id VARCHAR(36) NOT NULL,
    amount         BIGINT NOT NULL,       -- stored in cents
    status         VARCHAR(20) NOT NULL   -- Authorized, Declined
);
