package workers_storage

import (
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/dto"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/repositories/workers_repository"
	"time"
)

type WorkerStorage interface {
	Register(worker *dto.WorkerRequestDTO) (bool, error)
	FindAll() ([]*dto.WorkerResponseDTO, error)
	DeleteExpiredWorkers(deadline time.Time) ([]int, error)
	FindFreeWorker() (*dto.WorkerResponseDTO, error)
}

type workerStorage struct {
	repository workers_repository.WorkersRepository
}

func NewWorkerStorage(repository workers_repository.WorkersRepository) WorkerStorage {
	return &workerStorage{repository: repository}
}

func (w *workerStorage) Register(worker *dto.WorkerRequestDTO) (bool, error) {
	return w.repository.Register(&workers_repository.WorkerEntity{
		Id:        worker.Id,
		Url:       worker.Url,
		Executors: worker.Executors,
	})
}

func (w *workerStorage) FindAll() ([]*dto.WorkerResponseDTO, error) {
	entities, err := w.repository.FindAll()
	if err != nil {
		return nil, err
	}

	var workers []*dto.WorkerResponseDTO

	for _, e := range entities {
		workers = append(workers, &dto.WorkerResponseDTO{
			Id:           e.Id,
			Url:          e.Url,
			Executors:    e.Executors,
			LastModified: e.LastModified,
		})
	}

	if len(workers) == 0 {
		return []*dto.WorkerResponseDTO{}, nil
	}

	return workers, nil
}

func (w *workerStorage) DeleteExpiredWorkers(deadline time.Time) ([]int, error) {
	return w.repository.DeleteExpiredWorkers(deadline)
}

func (w *workerStorage) FindFreeWorker() (*dto.WorkerResponseDTO, error) {
	worker, err := w.repository.FindFreeWorker()
	if err != nil {
		return nil, err
	}

	if worker == nil {
		return nil, nil
	}

	if worker.Executors <= 0 {
		return nil, nil
	}

	return &dto.WorkerResponseDTO{
		Id:        worker.Id,
		Url:       worker.Url,
		Executors: worker.Executors,
	}, nil
}
