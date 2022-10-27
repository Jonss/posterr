CREATE TABLE users(
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(14) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT (now())
);

CREATE INDEX users_username_idx ON users (username);