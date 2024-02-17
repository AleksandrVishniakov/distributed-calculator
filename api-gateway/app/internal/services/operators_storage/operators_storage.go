package operators_storage

import (
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/dto"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/repositories/operator_repository"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/services/expression/expr_tokens"
)

type OperatorsStorage interface {
	SaveAll(operations []*dto.OperationDTO) error
	FindAll() ([]*dto.OperationDTO, error)
}

type operatorsStorage struct {
	repository operator_repository.OperatorsRepository
}

func NewOperatorsStorage(repository operator_repository.OperatorsRepository) OperatorsStorage {
	return &operatorsStorage{repository: repository}
}

func (o *operatorsStorage) SaveAll(operations []*dto.OperationDTO) error {
	for _, operation := range operations {
		err := o.repository.Save(&operator_repository.OperatorEntity{
			OperatorType: int(operation.OperationType),
			DurationMS:   operation.DurationMS,
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func (o *operatorsStorage) FindAll() ([]*dto.OperationDTO, error) {
	entities, err := o.repository.FindAll()
	if err != nil {
		return nil, err
	}

	var operations = []*dto.OperationDTO{}
	for _, entity := range entities {
		operations = append(operations, &dto.OperationDTO{
			OperationType: expr_tokens.OperationType(entity.OperatorType),
			DurationMS:    entity.DurationMS,
		})
	}

	return operations, nil
}
