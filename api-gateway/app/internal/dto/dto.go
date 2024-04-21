package dto

import (
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/services/expression/expr_tokens"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/services/statuses"
	"time"
)

type CalculationRequestDTO struct {
	Expression     string `json:"expression"`
	IdempotencyKey string `json:"idempotencyKey"`
}

type CalculationResponseDTO struct {
	Id int `json:"id"`
}

type ExpressionResponseDTO struct {
	Id         int       `json:"id"`
	Expression string    `json:"expression"`
	CreatedAt  time.Time `json:"createdAt"`
	FinishedAt time.Time `json:"finishedAt"`
	Status     int       `json:"status"`
	Result     float64   `json:"result"`
}

type OperationDTO struct {
	OperationType expr_tokens.OperationType `json:"operationType"`
	DurationMS    int                       `json:"durationMS"`
}

type WorkerRequestDTO struct {
	Id        uint64 `json:"id"`
	Url       string `json:"url"`
	Executors int    `json:"executors"`
}

type WorkerResponseDTO struct {
	Id           int       `json:"id"`
	Url          string    `json:"url"`
	Executors    int       `json:"executors"`
	LastModified time.Time `json:"lastModified"`
}

type CalculationResultDTO struct {
	Result float64 `json:"result"`
}

type ExpressionNodeDTO struct {
	Id            int             `json:"id"`
	UserID        uint64          `json:"userId"`
	ParentId      int             `json:"parentId"`
	ExpressionId  int             `json:"expressionId"`
	Type          int             `json:"type"`
	OperationType int             `json:"operationType"`
	Status        statuses.Status `json:"status"`
	Result        float64         `json:"result"`
	WorkerId      int             `json:"workerId"`
}

type TaskDTO struct {
	LeftResult    float64         `json:"leftResult"`
	OperationType int             `json:"operationType"`
	RightResult   float64         `json:"rightResult"`
	Status        statuses.Status `json:"status"`
}
