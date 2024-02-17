package executors_pool

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/AleksandrVishniakov/distributed-calculator/daemon/app/internal/dto"
	"github.com/AleksandrVishniakov/distributed-calculator/daemon/app/internal/services/operations"
	"log"
	"net/http"
	"sync"
	"time"
)

type CalculationExecutor struct {
	id               int
	first            float64
	second           float64
	operation        operations.OperationType
	duration         time.Duration
	orchestratorHost string
}

func NewCalculationExecutor(request *dto.CalculationRequestDTO, orchestratorHost string) *CalculationExecutor {
	return &CalculationExecutor{
		id:               request.Id,
		first:            request.First,
		second:           request.Second,
		operation:        request.Operation,
		duration:         request.Duration,
		orchestratorHost: orchestratorHost,
	}
}

func (e *CalculationExecutor) Task() {
	e.sendStartingRequest()

	var result float64

	switch e.operation {
	case operations.Plus:
		result = e.first + e.second
	case operations.Minus:
		result = e.first - e.second
	case operations.Multiply:
		result = e.first * e.second
	case operations.Divide:
		result = e.first / e.second
	}

	wg := sync.WaitGroup{}

	wg.Add(1)
	time.AfterFunc(e.duration, func() {
		defer wg.Done()
		e.sendResultRequest(result)
	})

	wg.Wait()
}

func (e *CalculationExecutor) sendStartingRequest() {
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	url := fmt.Sprintf("%s/api/task/%d/status", e.orchestratorHost, e.id)

	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode >= 400 {
		if err != nil {
			log.Fatal("starting request failed with error:", resp.Status)
		}
	}
}

func (e *CalculationExecutor) sendResultRequest(result float64) {
	requestJSON, err := json.Marshal(&dto.CalculationResultDTO{Result: result})
	if err != nil {
		log.Fatal(err)
	}

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	url := fmt.Sprintf("%s/api/task/%d/result", e.orchestratorHost, e.id)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestJSON))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode >= 400 {
		if err != nil {
			log.Fatal("starting request failed with error:", resp.Status)
		}
	}
}
