package dbhelper

import "database/sql"

func MapEntities[T any](rows *sql.Rows, mapper func(func(dest ...any) error) (*T, error)) (entities []*T, err error) {
	entities = []*T{}

	for rows.Next() {
		entity, err := mapper(rows.Scan)
		if err != nil {
			return nil, err
		}

		entities = append(entities, entity)
	}

	return entities, nil
}
