-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS expressions_tree (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    parent_id INT REFERENCES expressions_tree(id) ON DELETE CASCADE ON UPDATE CASCADE,
    expression_id INT REFERENCES expressions(id) ON DELETE CASCADE ON UPDATE CASCADE,
    type SMALLINT NOT NULL,
    operation_type SMALLINT,
    status SMALLINT NOT NULL,
    result double precision NOT NULL DEFAULT 0,
    worker_id INT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS expressions_tree;
-- +goose StatementEnd
