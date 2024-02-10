-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS workers (
    id SERIAL UNIQUE,
    url VARCHAR(1024) NOT NULL UNIQUE,
    executors INT NOT NULL,
    last_modified timestamptz NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS workers
-- +goose StatementEnd
