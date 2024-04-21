package grpcsrv

import (
	"context"
	orchestrator "github.com/AleksandrVishniakov/dc-protos/gen/go/orchestrator/v1"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/dto"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/services/binary_tree_storage"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/services/expressions_storage"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/services/operators_storage"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/services/worker_api"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/services/workers_storage"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/util/calc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const tempUserID = 1

type Server struct {
	orchestrator.UnimplementedOrchestratorServer

	binaryTreeStorage binary_tree_storage.BinaryTreeStorage
	operatorsStorage  operators_storage.OperatorsStorage
	workersStorage    workers_storage.WorkerStorage
	expressionStorage expressions_storage.ExpressionStorage

	workerAPI worker_api.WorkerAPI
}

func Register(
	gRPCServer *grpc.Server,

	binaryTreeStorage binary_tree_storage.BinaryTreeStorage,
	operatorsStorage operators_storage.OperatorsStorage,
	workersStorage workers_storage.WorkerStorage,
	expressionStorage expressions_storage.ExpressionStorage,

	workerAPI worker_api.WorkerAPI,
) {
	orchestrator.RegisterOrchestratorServer(gRPCServer, &Server{
		binaryTreeStorage: binaryTreeStorage,
		operatorsStorage:  operatorsStorage,
		workersStorage:    workersStorage,
		expressionStorage: expressionStorage,
		workerAPI:         workerAPI,
	})
}

func (s *Server) RegisterWorker(ctx context.Context, request *orchestrator.WorkerRegisterRequest) (*orchestrator.WorkerRegisterResponse, error) {
	exists, err := s.workersStorage.Register(&dto.WorkerRequestDTO{
		Id:        request.Id,
		Url:       request.Url,
		Executors: int(request.Executors),
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if exists {
		return &orchestrator.WorkerRegisterResponse{Ok: true}, nil
	}

	err = calc.CalculateAll(
		ctx,
		s.binaryTreeStorage,
		s.operatorsStorage,
		s.workersStorage,
		s.expressionStorage,
		s.workerAPI,
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &orchestrator.WorkerRegisterResponse{Ok: true}, nil
}

func (s *Server) StartTask(_ context.Context, request *orchestrator.TaskStartingRequest) (*orchestrator.TaskStartingResponse, error) {
	var id = int(request.GetId())

	err := s.binaryTreeStorage.MarkAsCalculating(id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	node, err := s.binaryTreeStorage.FindById(id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	err = s.expressionStorage.MarkAsCalculating(node.ExpressionId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &orchestrator.TaskStartingResponse{Ok: true}, nil
}

func (s *Server) SendTaskResult(ctx context.Context, request *orchestrator.TaskResultRequest) (*orchestrator.TaskResultResponse, error) {
	var id = int(request.GetId())

	err := s.binaryTreeStorage.SaveResult(id, float64(request.GetResult()))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	node, err := s.binaryTreeStorage.FindById(id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if node.ParentId == -1 {
		err = s.expressionStorage.SaveResult(node.ExpressionId, float64(request.GetResult()))
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

		return &orchestrator.TaskResultResponse{Ok: true}, nil
	}

	err = calc.CalculateAll(
		ctx,
		s.binaryTreeStorage,
		s.operatorsStorage,
		s.workersStorage,
		s.expressionStorage,
		s.workerAPI,
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &orchestrator.TaskResultResponse{Ok: true}, nil
}
