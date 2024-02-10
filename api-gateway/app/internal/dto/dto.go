package dto

import "time"

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

type LimitOffsetRequest struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type OperationsSetup struct {
	AddTime      int64 `json:"addTime"`
	SubtractTime int64 `json:"subtractTime"`
	MultiplyTime int64 `json:"multiplyTime"`
	DivideTime   int64 `json:"divideTime"`
}

type WorkerRequestDTO struct {
	Id        int    `json:"id"`
	Url       string `json:"url"`
	Executors int    `json:"executors"`
}

type WorkerResponseDTO struct {
	Id           int       `json:"id"`
	Url          string    `json:"url"`
	Executors    int       `json:"executors"`
	LastModified time.Time `json:"lastModified"`
}
