package worker_api

import (
	"context"
	"fmt"
	daemonv1 "github.com/AleksandrVishniakov/dc-protos/gen/go/daemon/v1"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/services/expression/expr_tokens"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type CalculationRequestDTO struct {
	Id        uint64                    `json:"id"`
	First     float64                   `json:"first"`
	Second    float64                   `json:"second"`
	Operation expr_tokens.OperationType `json:"operation"`
	Duration  time.Duration             `json:"duration"`
}

type WorkerAPI interface {
	Calculate(ctx context.Context, host string, userID uint64, requestBody *CalculationRequestDTO) error
}

type gRPCWorkerAPI struct{}

func NewGRPCWorkerAPI() WorkerAPI {
	return &gRPCWorkerAPI{}
}

func (g *gRPCWorkerAPI) Calculate(ctx context.Context, host string, userID uint64, requestBody *CalculationRequestDTO) error {
	cc, err := grpc.DialContext(
		ctx,
		host,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		return err
	}

	client := daemonv1.NewDaemonClient(cc)

	resp, err := client.CalculateTask(ctx, &daemonv1.CalculationRequestDTO{
		Id:        requestBody.Id,
		UserID:    userID,
		First:     float32(requestBody.First),
		Second:    float32(requestBody.Second),
		Operation: daemonv1.OperationType(requestBody.Operation),
		Duration:  uint64(requestBody.Duration.Milliseconds()),
	})
	if err != nil {
		return err
	}

	if !resp.Ok {
		return fmt.Errorf("response is not ok")
	}

	return nil
}
