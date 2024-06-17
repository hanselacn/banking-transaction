CREATE TABLE accounts (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    account_number VARCHAR(255) NOT NULL UNIQUE,
    balance VARCHAR(255) NOT NULL,
    interest_rate VARCHAR(255) NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
);