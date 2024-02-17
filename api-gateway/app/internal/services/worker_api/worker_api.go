package worker_api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/services/expression/expr_tokens"
	"net/http"
	"time"
)

type CalculationRequestDTO struct {
	Id        int                       `json:"id"`
	First     float64                   `json:"first"`
	Second    float64                   `json:"second"`
	Operation expr_tokens.OperationType `json:"operation"`
	Duration  time.Duration             `json:"duration"`
}

type WorkerAPI interface {
	Calculate(host string, requestBody *CalculationRequestDTO) error
}

type workerAPI struct {
	client http.Client
}

func NewWorkerAPI(timeout time.Duration) WorkerAPI {
	return &workerAPI{client: http.Client{Timeout: timeout}}
}

func (w *workerAPI) Calculate(host string, requestBody *CalculationRequestDTO) error {
	requestJSON, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/api/task", host)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestJSON))
	req.Header.Set("Content-Type", "application/json")

	resp, err := w.client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		return errors.New("worker calculate request failed with status: " + resp.Status)
	}

	return nil
}
