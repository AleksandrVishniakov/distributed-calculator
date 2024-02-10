package expressions_repository

import (
	"database/sql"
	"errors"
)

var (
	ErrExpressionNotFound = errors.New("expressions_repository: expression not found")
)

type ExpressionsRepository interface {
	Update(entity *ExpressionEntity) error
	FindAll() ([]*ExpressionEntity, error)
	FindById(id int) (*ExpressionEntity, error)
	FindByIdempotencyKey(key string, expression string) (int, error)
	Create(expressions string, status int, key string) (int, error)
}

type expressionsRepository struct {
	db *sql.DB
}

func NewExpressionsRepository(db *sql.DB) ExpressionsRepository {
	return &expressionsRepository{db: db}
}

func (e *expressionsRepository) FindByIdempotencyKey(key string, expression string) (int, error) {
	row := e.db.QueryRow(
		"SELECT id FROM expressions WHERE idempotency_key=$1 AND expression=$2 limit 1",
		key,
		expression,
	)

	var id int
	err := row.Scan(&id)

	if errors.Is(err, sql.ErrNoRows) {
		return 0, nil
	}

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (e *expressionsRepository) Create(expressions string, status int, key string) (int, error) {
	row := e.db.QueryRow(
		"INSERT INTO expressions (expression, status, idempotency_key) VALUES ($1, $2, $3) returning id",
		expressions,
		status,
		key,
	)

	var id int
	err := row.Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (e *expressionsRepository) FindById(id int) (*ExpressionEntity, error) {
	row := e.db.QueryRow(
		"SELECT * FROM expressions WHERE id=$1",
		id,
	)

	var entity = &ExpressionEntity{}
	err := row.Scan(
		&entity.Id,
		&entity.Expression,
		&entity.Status,
		&entity.Result,
		&entity.CreatedAt,
		&entity.FinishedAt,
		&entity.IdempotencyKey,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrExpressionNotFound
	}

	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (e *expressionsRepository) Update(entity *ExpressionEntity) error {
	_, err := e.db.Exec(
		"UPDATE expressions SET expression=$1, status=$2, result=$3, finished_at=$4, idempotency_key=$5 WHERE id=$6",
		entity.Expression,
		entity.Status,
		entity.Result,
		entity.FinishedAt,
		entity.IdempotencyKey,
		entity.Id,
	)

	return err
}

func (e *expressionsRepository) FindAll() ([]*ExpressionEntity, error) {
	rows, err := e.db.Query("SELECT * FROM expressions")
	if err != nil {
		return []*ExpressionEntity{}, err
	}

	var expressions []*ExpressionEntity

	for rows.Next() {
		var expr = &ExpressionEntity{}

		err := rows.Scan(
			&expr.Id,
			&expr.Expression,
			&expr.Status,
			&expr.Result,
			&expr.CreatedAt,
			&expr.FinishedAt,
			&expr.IdempotencyKey,
		)

		if err != nil {
			return []*ExpressionEntity{}, err
		}

		expressions = append(expressions, expr)
	}

	return expressions, nil
}
