package expressions_storage

import (
	"errors"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/dto"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/repositories/expressions_repository"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/services/expression"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/services/statuses"
	"time"
)

var (
	ErrExpressionNotFound = errors.New("expressions_storage: expression not found")
)

type ExpressionStorage interface {
	FindByIdempotencyKey(key string, expression expression.Expression) (int, error)
	Create(expressions expression.Expression, userID uint64, key string) (int, error)
	FindById(id int) (*dto.ExpressionResponseDTO, error)
	FindAll() ([]*dto.ExpressionResponseDTO, error)
	FindAllByUserID(userID uint64) ([]*dto.ExpressionResponseDTO, error)
	SaveResult(id int, result float64) error
	MarkAsCalculating(id int) error
	MarkAsFailed(id int) error
}

type expressionStorage struct {
	repository expressions_repository.ExpressionsRepository
}

func NewExpressionStorage(repository expressions_repository.ExpressionsRepository) ExpressionStorage {
	return &expressionStorage{repository: repository}
}

func (e *expressionStorage) FindByIdempotencyKey(key string, expression expression.Expression) (int, error) {
	return e.repository.FindByIdempotencyKey(key, string(expression))
}

func (e *expressionStorage) Create(expr expression.Expression, userID uint64, key string) (int, error) {
	return e.repository.Create(string(expr), userID, int(statuses.Created), key)
}

func (e *expressionStorage) FindById(id int) (*dto.ExpressionResponseDTO, error) {
	entity, err := e.repository.FindById(id)

	if errors.Is(err, expressions_repository.ErrExpressionNotFound) {
		return nil, ErrExpressionNotFound
	}

	if err != nil {
		return nil, err
	}

	return dto.MapExpressionResponseFromEntity(entity), nil
}

func (e *expressionStorage) FindAll() ([]*dto.ExpressionResponseDTO, error) {
	entities, err := e.repository.FindAll()
	if err != nil {
		return nil, err
	}

	var expressions []*dto.ExpressionResponseDTO

	for _, e := range entities {
		expressions = append(expressions, dto.MapExpressionResponseFromEntity(e))
	}

	return expressions, nil
}

func (e *expressionStorage) FindAllByUserID(userID uint64) ([]*dto.ExpressionResponseDTO, error) {
	entities, err := e.repository.FindAllByUserID(userID)
	if err != nil {
		return nil, err
	}

	var expressions []*dto.ExpressionResponseDTO

	for _, e := range entities {
		expressions = append(expressions, dto.MapExpressionResponseFromEntity(e))
	}

	return expressions, nil
}

func (e *expressionStorage) SaveResult(id int, result float64) error {
	return e.repository.Update(&expressions_repository.ExpressionEntity{
		Id:         id,
		Result:     result,
		Status:     int(statuses.Finished),
		FinishedAt: time.Now(),
	})
}

func (e *expressionStorage) MarkAsCalculating(id int) error {
	return e.repository.SetStatus(id, int(statuses.Calculating))
}

func (e *expressionStorage) MarkAsFailed(id int) error {
	return e.repository.SetStatus(id, int(statuses.Failed))
}
