package dto

import (
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/repositories/expressions_repository"
)

func MapExpressionResponseFromEntity(entity *expressions_repository.ExpressionEntity) *ExpressionResponseDTO {
	return &ExpressionResponseDTO{
		Id:         entity.Id,
		Expression: entity.Expression,
		CreatedAt:  entity.CreatedAt,
		FinishedAt: entity.FinishedAt,
		Status:     entity.Status,
		Result:     entity.Result,
	}
}
