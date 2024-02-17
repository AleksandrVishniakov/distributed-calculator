package workers_repository

import (
	"database/sql"
	"errors"
	"time"
)

type WorkersRepository interface {
	Register(entity *WorkerEntity) (bool, error)
	FindAll() ([]*WorkerEntity, error)
	DeleteExpiredWorkers(deadline time.Time) ([]int, error)
	FindFreeWorker() (*FreeWorkerEntity, error)
}

type workersRepository struct {
	db *sql.DB
}

func NewWorkersRepository(db *sql.DB) WorkersRepository {
	return &workersRepository{db: db}
}

func (w *workersRepository) Register(entity *WorkerEntity) (bool, error) {
	row := w.db.QueryRow(
		"INSERT INTO workers (id, url, executors) VALUES ($1, $2, $3) ON CONFLICT (id) DO UPDATE SET url = $2, executors = $3, last_modified = NOW() returning xmax::text::int > 0 as is_updated",
		entity.Id,
		entity.Url,
		entity.Executors,
	)

	var exists bool

	err := row.Scan(&exists)

	return exists, err
}

func (w *workersRepository) FindAll() ([]*WorkerEntity, error) {
	rows, err := w.db.Query("SELECT * FROM workers ORDER BY id")
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

func (w *workersRepository) FindFreeWorker() (*FreeWorkerEntity, error) {
	row := w.db.QueryRow(
		`select
		w.id, w.url, w.executors - (select count(*) from expressions_tree t where t.worker_id = w.id and t.status <> 3 and t.status <> 4) as free_executors
		from workers w order by free_executors desc, w.id asc limit 1`,
	)

	var worker = &FreeWorkerEntity{}

	err := row.Scan(&worker.Id, &worker.Url, &worker.Executors)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return worker, nil
}
