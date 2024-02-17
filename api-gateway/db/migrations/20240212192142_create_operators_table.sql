-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS operators (
    operator_type int NOT NULL UNIQUE,
    duration_ms int NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS operators;
-- +goose StatementEnd
