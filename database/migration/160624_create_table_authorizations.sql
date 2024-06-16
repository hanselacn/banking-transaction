CREATE TABLE authorizations (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    password VARCHAR(255) NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
);