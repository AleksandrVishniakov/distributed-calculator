package operator_repository

import "database/sql"

type OperatorsRepository interface {
	Save(entity *OperatorEntity) error
	FindAll() ([]*OperatorEntity, error)
}

type operatorsRepository struct {
	db *sql.DB
}

func NewOperatorsRepository(db *sql.DB) OperatorsRepository {
	return &operatorsRepository{
		db: db,
	}
}

func (o *operatorsRepository) Save(entity *OperatorEntity) error {
	_, err := o.db.Exec(
		"INSERT INTO operators (operator_type, duration_ms) VALUES ($1, $2) ON CONFLICT (operator_type) DO UPDATE SET duration_ms = $2",
		entity.OperatorType,
		entity.DurationMS,
	)

	return err
}

func (o *operatorsRepository) FindAll() ([]*OperatorEntity, error) {
	rows, err := o.db.Query(
		"SELECT * FROM operators",
	)

	if err != nil {
		return nil, err
	}

	var operators = []*OperatorEntity{}

	for rows.Next() {
		var operator = &OperatorEntity{}

		err := rows.Scan(&operator.OperatorType, &operator.DurationMS)
		if err != nil {
			return nil, err
		}

		operators = append(operators, operator)
	}

	return operators, nil
}
