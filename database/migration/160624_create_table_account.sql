CREATE TABLE accounts (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    account_number VARCHAR(255) NOT NULL UNIQUE,
    balance NUMERIC(15, 2) NOT NULL,
    interest_rate NUMERIC(15, 2) NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
);