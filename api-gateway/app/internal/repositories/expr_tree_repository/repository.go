package expr_tree_repository

import (
	"database/sql"
)

type ExpressionsTreeRepository interface {
	Create(entity *ExpressionTreeNodeEntity) (int, error)
	SetStatus(id int, status int) error
	SaveResult(id int, result float64, status int) error
	FindByParentId(parentId int) ([]*ExpressionTreeNodeEntity, error)
	SaveWorker(id int, workerId int, status int) error
	FindById(id int) (*ExpressionTreeNodeEntity, error)
	FindByWorkerId(id int) ([]*TaskEntity, error)
	DeleteWorker(workerId int) error
	DeleteAllWorkers() error
	FindUncalculated() ([]int, error)
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

func (e *expressionsTreeRepository) SetStatus(id int, status int) error {
	_, err := e.db.Exec(
		"UPDATE expressions_tree SET status=$1 WHERE id=$2 and status < $1",
		status,
		id,
	)

	return err
}

func (e *expressionsTreeRepository) SaveResult(id int, result float64, status int) error {
	_, err := e.db.Exec(
		"UPDATE expressions_tree SET result=$1, status=$2 WHERE id=$3",
		result,
		status,
		id,
	)

	return err
}

func nullableInt(n int) sql.NullInt32 {
	if n != -1 {
		return sql.NullInt32{Int32: int32(n), Valid: true}
	}
	return sql.NullInt32{}
}

func (e *expressionsTreeRepository) FindByParentId(parentId int) ([]*ExpressionTreeNodeEntity, error) {
	rows, err := e.db.Query(
		"SELECT * FROM expressions_tree WHERE parent_id = $1 ORDER BY type",
		parentId,
	)

	if err != nil {
		return nil, err
	}

	var entities []*ExpressionTreeNodeEntity

	for rows.Next() {
		var entity = &ExpressionTreeNodeEntity{}

		var nullableParentId sql.NullInt32
		var nullableWorkerId sql.NullInt32
		var nullableOperationType sql.NullInt32

		err := rows.Scan(&entity.Id, &nullableParentId, &entity.ExpressionId, &entity.Type, &nullableOperationType, &entity.Status, &entity.Result, &nullableWorkerId)
		if err != nil {
			return nil, err
		}

		entity.ParentId = int(nullableParentId.Int32)
		entity.WorkerId = int(nullableWorkerId.Int32)

		if nullableOperationType.Valid {
			entity.OperationType = int(nullableOperationType.Int32)
		} else {
			entity.OperationType = -1
		}

		entities = append(entities, entity)
	}

	return entities, nil
}

func (e *expressionsTreeRepository) SaveWorker(id int, workerId int, status int) error {
	_, err := e.db.Exec(
		"UPDATE expressions_tree SET status=$1, worker_id=$2 WHERE id=$3",
		status,
		workerId,
		id,
	)

	return err
}

func (e *expressionsTreeRepository) FindById(id int) (*ExpressionTreeNodeEntity, error) {
	row := e.db.QueryRow(
		"SELECT * FROM expressions_tree WHERE id = $1",
		id,
	)

	var entity = &ExpressionTreeNodeEntity{}

	var nullableParentId sql.NullInt32
	var nullableWorkerId sql.NullInt32
	var nullableOperationType sql.NullInt32

	err := row.Scan(&entity.Id, &nullableParentId, &entity.ExpressionId, &entity.Type, &nullableOperationType, &entity.Status, &entity.Result, &nullableWorkerId)
	entity.WorkerId = int(nullableWorkerId.Int32)

	if nullableParentId.Valid {
		entity.ParentId = int(nullableParentId.Int32)
	} else {
		entity.ParentId = -1
	}

	if nullableOperationType.Valid {
		entity.OperationType = int(nullableOperationType.Int32)
	} else {
		entity.OperationType = -1
	}

	return entity, err
}

func (e *expressionsTreeRepository) FindByWorkerId(workerId int) ([]*TaskEntity, error) {
	rows, err := e.db.Query(
		`select 
    			(select l.result from expressions_tree l where l.parent_id = op.id and l.type = 0) as left_result, 
    			op.operation_type, 
    			(select r.result from expressions_tree r where r.parent_id = op.id and r.type = 1) as right_result, 
    			op.status 
				from expressions_tree op where op.worker_id = $1 and op.status <> 3 and op.status <> 4`,
		workerId,
	)

	if err != nil {
		return nil, err
	}

	var entities []*TaskEntity

	for rows.Next() {
		var entity = &TaskEntity{}

		err := rows.Scan(&entity.LeftResult, &entity.OperationType, &entity.RightResult, &entity.Status)
		if err != nil {
			return nil, err
		}

		entities = append(entities, entity)
	}

	return entities, nil
}

func (e *expressionsTreeRepository) DeleteWorker(workerId int) error {
	_, err := e.db.Exec(
		"UPDATE expressions_tree SET worker_id = null, status=0 WHERE worker_id = $1 AND status <> 3 AND status <> 4",
		workerId,
	)

	return err
}

func (e *expressionsTreeRepository) DeleteAllWorkers() error {
	_, err := e.db.Exec(
		"UPDATE expressions_tree SET worker_id = null, status=0 WHERE status <> 3 AND status <> 4 AND worker_id IS NOT NULL",
	)

	return err
}

func (e *expressionsTreeRepository) FindUncalculated() ([]int, error) {
	rows, err := e.db.Query(
		`select id from expressions_tree where parent_id IS NULL AND status = 0`,
	)

	if err != nil {
		return []int{}, err
	}

	var ids []int

	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			return []int{}, err
		}

		ids = append(ids, id)
	}

	return ids, nil
}
