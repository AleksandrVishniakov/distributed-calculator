package executors_pool

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"sync"
	"time"

	orchestrator "github.com/AleksandrVishniakov/dc-protos/gen/go/orchestrator/v1"
	"github.com/AleksandrVishniakov/distributed-calculator/daemon/app/internal/dto"
	"github.com/AleksandrVishniakov/distributed-calculator/daemon/app/internal/services/operations"
)

type CalculationExecutor struct {
	id        uint64
	first     float64
	second    float64
	operation operations.OperationType
	duration  time.Duration

	client orchestrator.OrchestratorClient
}

func NewCalculationExecutor(
	ctx context.Context,
	request *dto.CalculationRequestDTO,
	orchestratorGRPCHost string,
) (*CalculationExecutor, error) {
	cc, err := grpc.DialContext(
		ctx,
		orchestratorGRPCHost,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return &CalculationExecutor{
		id:        request.ID,
		first:     request.First,
		second:    request.Second,
		operation: request.Operation,
		duration:  request.Duration,

		client: orchestrator.NewOrchestratorClient(cc),
	}, nil
}

func (e *CalculationExecutor) Task(ctx context.Context) {
	e.sendStartingRequest(ctx)

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
		e.sendResultRequest(ctx, result)
	})

	wg.Wait()
}

func (e *CalculationExecutor) sendStartingRequest(ctx context.Context) {
	log.Println("starting request")
	resp, err := e.client.StartTask(ctx, &orchestrator.TaskStartingRequest{
		Id: e.id,
	})

	if err != nil {
		log.Fatal("starting request err: ", err)
	}

	if !resp.Ok {
		log.Fatal("starting request is not ok")
	}
}

func (e *CalculationExecutor) sendResultRequest(ctx context.Context, result float64) {
	resp, err := e.client.SendTaskResult(ctx, &orchestrator.TaskResultRequest{
		Id:     e.id,
		Result: float32(result),
	})

	if err != nil {
		log.Fatal("starting request err: ", err)
	}

	if !resp.Ok {
		log.Fatal("starting request is not ok")
	}
}
