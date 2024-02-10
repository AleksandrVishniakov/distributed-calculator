package binary_tree_storage

import (
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/repositories/expr_tree_repository"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/services/binary_tree"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/services/expression/expr_tokens"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/services/statuses"
)

type BinaryTreeStorage interface {
	SaveTree(root *binary_tree.Node, expressionId int, parentId int) error
}

type binaryTreeStorage struct {
	repository expr_tree_repository.ExpressionsTreeRepository
}

func NewBinaryTreeStorage(repository expr_tree_repository.ExpressionsTreeRepository) BinaryTreeStorage {
	return &binaryTreeStorage{repository: repository}
}

func (b *binaryTreeStorage) SaveTree(node *binary_tree.Node, expressionId int, parentId int) error {
	if node == nil {
		return nil
	}

	var status int
	switch node.Value.Type() {
	case expr_tokens.BinaryOperation:
		status = int(statuses.InProgress)
	default:
		status = int(statuses.Finished)
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
		Type:          int(node.Value.Type()),
		OperationType: operationType,
		Status:        status,
		Result:        result,
	})

	if err != nil {
		return err
	}

	err = b.SaveTree(node.Left, expressionId, id)
	if err != nil {
		return err
	}

	err = b.SaveTree(node.Right, expressionId, id)
	if err != nil {
		return err
	}

	return nil
}
