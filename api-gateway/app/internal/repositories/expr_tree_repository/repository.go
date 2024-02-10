package expr_tree_repository

import (
	"database/sql"
)

type ExpressionsTreeRepository interface {
	Create(entity *ExpressionTreeNodeEntity) (int, error)
}

type expressionsTreeRepository struct {
	db *sql.DB
}

func NewExpressionsTreeRepository(db *sql.DB) ExpressionsTreeRepository {
	return &expressionsTreeRepository{db: db}
}

func (e *expressionsTreeRepository) Create(entity *ExpressionTreeNodeEntity) (int, error) {
	row := e.db.QueryRow(
		"INSERT INTO expressions_tree (parent_id, expression_id, type, operation_type, status, result) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		nullableInt(entity.ParentId),
		entity.ExpressionId,
		entity.Type,
		nullableInt(entity.OperationType),
		entity.Status,
		entity.Result,
	)

	var id int
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func nullableInt(n int) sql.NullInt32 {
	if n != -1 {
		return sql.NullInt32{Int32: int32(n), Valid: true}
	}
	return sql.NullInt32{}
}
