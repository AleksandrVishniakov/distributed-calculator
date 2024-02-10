package workers_repository

import "time"

type WorkerEntity struct {
	Id           int
	Url          string
	Executors    int
	LastModified time.Time
}
