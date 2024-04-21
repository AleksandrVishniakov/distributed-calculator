package calc

import (
	"context"
	"errors"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/dto"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/services/binary_tree_storage"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/services/expression/expr_tokens"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/services/expressions_storage"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/services/operators_storage"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/services/statuses"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/services/worker_api"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/services/workers_storage"
	"time"
)

func StartCalculating(
	ctx context.Context,
	taskID int,
	binaryTreeStorage binary_tree_storage.BinaryTreeStorage,
	operatorsStorage operators_storage.OperatorsStorage,
	workersStorage workers_storage.WorkerStorage,
	expressionStorage expressions_storage.ExpressionStorage,

	workerAPI worker_api.WorkerAPI,
) error {
	node, err := binaryTreeStorage.FindById(taskID)
	if err != nil {
		return err
	}

	var userID = node.UserID

	if node.Status != statuses.Created {
		return nil
	}

	nodes, err := binaryTreeStorage.FindByParentId(taskID)
	if err != nil {
		return err
	}

	if len(nodes) == 0 {
		return nil
	}

	left := nodes[0]
	right := nodes[1]

	if left.Status == statuses.Finished && right.Status == statuses.Finished {
		operations, err := operatorsStorage.FindAll()
		if err != nil {
			return err
		}

		operationDuration, err := getOperationDuration(operations, expr_tokens.OperationType(node.OperationType))
		if err != nil {
			return err
		}

		worker, err := workersStorage.FindFreeWorker()
		if err != nil {
			return err
		}

		if worker == nil {
			return nil
		}

		var operation = expr_tokens.OperationType(node.OperationType)

		if operation == expr_tokens.Divide && right.Result == 0 {
			err = expressionStorage.MarkAsFailed(node.ExpressionId)
			if err != nil {
				return err
			}

			err = binaryTreeStorage.MarkAsFailed(node.Id)
			if err != nil {
				return err
			}

			return nil
		}

		err = workerAPI.Calculate(ctx, worker.Url, userID, &worker_api.CalculationRequestDTO{
			Id:        uint64(taskID),
			First:     left.Result,
			Second:    right.Result,
			Operation: operation,
			Duration:  operationDuration,
		})
		if err != nil {
			return err
		}

		err = binaryTreeStorage.SaveWorker(taskID, worker.Id)
		if err != nil {
			return err
		}
	} else {
		if left.Status == statuses.Created {
			err := StartCalculating(ctx, left.Id, binaryTreeStorage, operatorsStorage, workersStorage, expressionStorage, workerAPI)
			if err != nil {
				return err
			}
		}

		if right.Status == statuses.Created {
			err := StartCalculating(ctx, right.Id, binaryTreeStorage, operatorsStorage, workersStorage, expressionStorage, workerAPI)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func CalculateAll(
	ctx context.Context,
	binaryTreeStorage binary_tree_storage.BinaryTreeStorage,
	operatorsStorage operators_storage.OperatorsStorage,
	workersStorage workers_storage.WorkerStorage,
	expressionStorage expressions_storage.ExpressionStorage,

	workerAPI worker_api.WorkerAPI,
) error {
	ids, err := binaryTreeStorage.FindUncalculated()
	if err != nil {
		return err
	}

	for _, id := range ids {
		err = StartCalculating(ctx, id, binaryTreeStorage, operatorsStorage, workersStorage, expressionStorage, workerAPI)
		if err != nil {
			return err
		}
	}

	return nil
}

func getOperationDuration(operations []*dto.OperationDTO, operationType expr_tokens.OperationType) (time.Duration, error) {
	for _, operation := range operations {
		if operation.OperationType == operationType {
			return time.Duration(operation.DurationMS) * time.Millisecond, nil
		}
	}

	return 0, errors.New("operation not found")
}
