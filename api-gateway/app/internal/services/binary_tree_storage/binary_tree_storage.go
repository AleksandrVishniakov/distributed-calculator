package binary_tree_storage

import (
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/dto"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/repositories/expr_tree_repository"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/services/binary_tree"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/services/expression/expr_tokens"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/services/statuses"
)

type BinaryTreeStorage interface {
	SaveTree(root *binary_tree.Node, expressionId int, parentId int, isLeft bool) (int, error)
	MarkAsCalculating(id int) error
	MarkAsFailed(id int) error
	SaveResult(id int, result float64) error
	FindByParentId(parentId int) ([]*dto.ExpressionNodeDTO, error)
	SaveWorker(id int, workerId int) error
	FindById(id int) (*dto.ExpressionNodeDTO, error)
	FindByWorkerId(id int) ([]*dto.TaskDTO, error)
	DeleteWorkers(workerIds []int) error
	DeleteAllWorkers() error
	FindUncalculated() ([]int, error)
}

type binaryTreeStorage struct {
	repository expr_tree_repository.ExpressionsTreeRepository
}

func NewBinaryTreeStorage(repository expr_tree_repository.ExpressionsTreeRepository) BinaryTreeStorage {
	return &binaryTreeStorage{repository: repository}
}

func (b *binaryTreeStorage) SaveTree(node *binary_tree.Node, expressionId int, parentId int, isLeft bool) (int, error) {
	if node == nil {
		return 0, nil
	}

	var status int
	switch node.Value.Type() {
	case expr_tokens.BinaryOperation:
		status = int(statuses.Created)
	default:
		status = int(statuses.Finished)
	}

	var taskType int
	if isLeft {
		taskType = 0
	} else {
		taskType = 1
	}

	var result float64
	if node.Value.Type() == expr_tokens.Number {
		result = node.Value.(*expr_tokens.NumberToken).Value
	}

	var operationType = -1
	if node.Value.Type() == expr_tokens.BinaryOperation {
		operationType = int(node.Value.(*expr_tokens.BinaryOperationToken).Operation)
	}

	id, err := b.repository.Create(&expr_tree_repository.ExpressionTreeNodeEntity{
		ParentId:      parentId,
		ExpressionId:  expressionId,
		Type:          taskType,
		OperationType: operationType,
		Status:        status,
		Result:        result,
	})

	if err != nil {
		return 0, err
	}

	_, err = b.SaveTree(node.Left, expressionId, id, true)
	if err != nil {
		return 0, err
	}

	_, err = b.SaveTree(node.Right, expressionId, id, false)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (b *binaryTreeStorage) MarkAsCalculating(id int) error {
	return b.repository.SetStatus(id, int(statuses.Calculating))
}

func (b *binaryTreeStorage) MarkAsFailed(id int) error {
	return b.repository.SetStatus(id, int(statuses.Failed))
}

func (b *binaryTreeStorage) SaveResult(id int, result float64) error {
	return b.repository.SaveResult(id, result, int(statuses.Finished))
}

func (b *binaryTreeStorage) FindByParentId(parentId int) ([]*dto.ExpressionNodeDTO, error) {
	entities, err := b.repository.FindByParentId(parentId)
	if err != nil {
		return nil, err
	}

	var expressionNodes []*dto.ExpressionNodeDTO

	for _, entity := range entities {
		expressionNodes = append(expressionNodes, &dto.ExpressionNodeDTO{
			Id:            entity.Id,
			ParentId:      entity.ParentId,
			ExpressionId:  entity.ExpressionId,
			Type:          entity.Type,
			OperationType: entity.OperationType,
			Status:        statuses.Status(entity.Status),
			Result:        entity.Result,
			WorkerId:      entity.WorkerId,
		})
	}

	return expressionNodes, nil
}

func (b *binaryTreeStorage) SaveWorker(id int, workerId int) error {
	return b.repository.SaveWorker(id, workerId, int(statuses.Enqueued))
}

func (b *binaryTreeStorage) FindById(id int) (*dto.ExpressionNodeDTO, error) {
	entity, err := b.repository.FindById(id)

	return &dto.ExpressionNodeDTO{
		Id:            entity.Id,
		ParentId:      entity.ParentId,
		ExpressionId:  entity.ExpressionId,
		Type:          entity.Type,
		OperationType: entity.OperationType,
		Status:        statuses.Status(entity.Status),
		Result:        entity.Result,
		WorkerId:      entity.WorkerId,
	}, err
}

func (b *binaryTreeStorage) FindByWorkerId(workerId int) ([]*dto.TaskDTO, error) {
	entities, err := b.repository.FindByWorkerId(workerId)
	if err != nil {
		return nil, err
	}

	var tasks = []*dto.TaskDTO{}

	for _, entity := range entities {
		tasks = append(tasks, &dto.TaskDTO{
			LeftResult:    entity.LeftResult,
			OperationType: entity.OperationType,
			RightResult:   entity.RightResult,
			Status:        statuses.Status(entity.Status),
		})
	}

	return tasks, nil
}

func (b *binaryTreeStorage) DeleteWorkers(workerIds []int) error {
	for _, workerId := range workerIds {
		err := b.repository.DeleteWorker(workerId)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *binaryTreeStorage) DeleteAllWorkers() error {
	return b.repository.DeleteAllWorkers()
}

func (b *binaryTreeStorage) FindUncalculated() ([]int, error) {
	return b.repository.FindUncalculated()
}
