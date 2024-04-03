-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    login VARCHAR(128) UNIQUE,
    password_hash TEXT UNIQUE
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_login
ON users(login);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
