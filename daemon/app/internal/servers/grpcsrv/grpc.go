package grpcsrv

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"

	daemonsrv "github.com/AleksandrVishniakov/dc-protos/gen/go/daemon/v1"
	dtos "github.com/AleksandrVishniakov/distributed-calculator/daemon/app/internal/dto"
	"github.com/AleksandrVishniakov/distributed-calculator/daemon/app/internal/services/executors_pool"
	"github.com/AleksandrVishniakov/distributed-calculator/daemon/app/internal/services/operations"
	"google.golang.org/grpc"
)

type Server struct {
	daemonsrv.UnimplementedDaemonServer

	orchestratorGRPCHost string
	poolManager          *executors_pool.PoolManager
}

func Register(
	gRPCServer *grpc.Server,
	orchestratorHost string,
	poolManager *executors_pool.PoolManager,
) {
	daemonsrv.RegisterDaemonServer(gRPCServer, &Server{
		orchestratorGRPCHost: orchestratorHost,
		poolManager:          poolManager,
	})
}

func (s *Server) CalculateTask(ctx context.Context, dto *daemonsrv.CalculationRequestDTO) (*daemonsrv.CalculationResponseDTO, error) {
	pool := s.poolManager.Pool(dto.UserID)

	executor, err := executors_pool.NewCalculationExecutor(ctx, &dtos.CalculationRequestDTO{
		ID:        dto.Id,
		UserID:    dto.UserID,
		First:     float64(dto.First),
		Second:    float64(dto.Second),
		Operation: operations.OperationType(dto.Operation),
		Duration:  time.Duration(dto.Duration) * time.Millisecond,
	}, s.orchestratorGRPCHost)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	pool.Run(executor)

	return &daemonsrv.CalculationResponseDTO{Ok: true}, nil
}
