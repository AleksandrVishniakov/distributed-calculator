package expressions_repository

import "time"

type ExpressionEntity struct {
	Id             int
	Expression     string
	Result         float64
	Status         int
	CreatedAt      time.Time
	FinishedAt     time.Time
	IdempotencyKey string
}
