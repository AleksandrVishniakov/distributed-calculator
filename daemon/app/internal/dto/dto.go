package dto

import (
	"github.com/AleksandrVishniakov/distributed-calculator/daemon/app/internal/services/operations"
	"time"
)

type CalculationRequestDTO struct {
	ID        uint64                   `json:"id"`
	UserID    uint64                   `json:"userID"`
	First     float64                  `json:"first"`
	Second    float64                  `json:"second"`
	Operation operations.OperationType `json:"operation"`
	Duration  time.Duration            `json:"duration"`
}

type OrchestratorPingDTO struct {
	ID        uint64 `json:"id"`
	Url       string `json:"url"`
	Executors int    `json:"executors"`
}

type CalculationResultDTO struct {
	Result float64 `json:"result"`
}
