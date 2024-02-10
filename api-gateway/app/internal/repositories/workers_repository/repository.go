package workers_repository

import (
	"database/sql"
	"time"
)

type WorkersRepository interface {
	Register(entity *WorkerEntity) error
	FindAll() ([]*WorkerEntity, error)
	DeleteExpiredWorkers(deadline time.Time) ([]int, error)
}

type workersRepository struct {
	db *sql.DB
}

func NewWorkersRepository(db *sql.DB) WorkersRepository {
	return &workersRepository{db: db}
}

func (w *workersRepository) Register(entity *WorkerEntity) error {
	_, err := w.db.Exec(
		"INSERT INTO workers (id, url, executors) VALUES ($1, $2, $3) ON CONFLICT (id) DO UPDATE SET url = $2, executors = $3, last_modified = NOW()",
		entity.Id,
		entity.Url,
		entity.Executors,
	)

	return err
}

func (w *workersRepository) FindAll() ([]*WorkerEntity, error) {
	rows, err := w.db.Query("SELECT * FROM workers")
	if err != nil {
		return []*WorkerEntity{}, err
	}

	var workers []*WorkerEntity

	for rows.Next() {
		var worker = &WorkerEntity{}

		err := rows.Scan(
			&worker.Id,
			&worker.Url,
			&worker.Executors,
			&worker.LastModified,
		)

		if err != nil {
			return []*WorkerEntity{}, err
		}

		workers = append(workers, worker)
	}

	return workers, nil
}

func (w *workersRepository) DeleteExpiredWorkers(deadline time.Time) ([]int, error) {
	rows, err := w.db.Query(
		"DELETE FROM workers WHERE last_modified < $1 returning id",
		deadline,
	)

	if err != nil {
		return nil, err
	}

	var ids []int

	for rows.Next() {
		var id int

		err := rows.Scan(&id)
		if err != nil {
			return nil, err
		}

		ids = append(ids, id)
	}

	return ids, nil
}
