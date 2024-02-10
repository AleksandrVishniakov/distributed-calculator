-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS expressions (
    id SERIAL PRIMARY KEY,
    expression TEXT NOT NULL DEFAULT '',
    status smallint NOT NULL DEFAULT 0,
    result double precision NOT NULL DEFAULT 0,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    finished_at timestamptz NOT NULL DEFAULT NOW(),
    idempotency_key VARCHAR(36) DEFAULT ''
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS expressions;
-- +goose StatementEnd
