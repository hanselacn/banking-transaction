CREATE TYPE transaction_type AS ENUM ('D', 'C');
CREATE TYPE transaction_action AS ENUM ('WITHDRAWAL', 'DEPOSIT','TRANSFER','PURCHASE');
CREATE TYPE transaction_status AS ENUM ('IN_PROGRESS','COMPLETED','FAILED');

CREATE TABLE transactions (
    id UUID PRIMARY KEY,
    type transaction_type,
    amount NUMERIC(15, 2) NOT NULL,
    action transaction_action,
    status transaction_status,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);