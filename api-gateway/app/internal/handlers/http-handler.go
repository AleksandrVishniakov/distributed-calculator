package handlers

import (
	"errors"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/dto"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/services/binary_tree"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/services/binary_tree_storage"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/services/expression"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/services/expressions_storage"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/services/workers_storage"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type HTTPHandler struct {
	expressionStorage expressions_storage.ExpressionStorage
	binaryTreeStorage binary_tree_storage.BinaryTreeStorage
	workersStorage    workers_storage.WorkerStorage
}

func NewHTTPHandler(
	expressionStorage expressions_storage.ExpressionStorage,
	binaryTreeStorage binary_tree_storage.BinaryTreeStorage,
	workersStorage workers_storage.WorkerStorage,
) *HTTPHandler {
	return &HTTPHandler{
		expressionStorage: expressionStorage,
		binaryTreeStorage: binaryTreeStorage,
		workersStorage:    workersStorage,
	}
}

func (h *HTTPHandler) InitRoutes() *gin.Engine {
	router := gin.Default()

	api := router.Group("/api")
	{
		api.POST("calculate", h.calculateExpression)
		tasks := api.Group("/tasks")
		{
			tasks.GET("/", h.getAllExpressions)
			tasks.GET("/status/:id", h.handleExpressionStatusRequest)
		}

		api.POST("/worker", h.handleWorkerRegister)
		api.GET("/workers", h.getAllWorkers)
	}

	return router
}

func (h *HTTPHandler) calculateExpression(c *gin.Context) {
	var calculationRequest dto.CalculationRequestDTO

	err := c.BindJSON(&calculationRequest)
	if err != nil {
		dto.NewResponseError(http.StatusBadRequest, err.Error()).Abort(c)
		return
	}

	if calculationRequest.Expression == "" {
		dto.NewResponseError(http.StatusBadRequest, "expression is empty").Abort(c)
		return
	}

	expr, err := expression.NewExpression(calculationRequest.Expression)
	if err != nil {
		dto.NewResponseError(http.StatusBadRequest, err.Error()).Abort(c)
		return
	}

	if calculationRequest.IdempotencyKey != "" {
		id, err := h.expressionStorage.FindByIdempotencyKey(calculationRequest.IdempotencyKey, expr)
		if err != nil {
			dto.NewResponseError(http.StatusInternalServerError, err.Error()).Abort(c)
			return
		}

		if id != 0 {
			c.IndentedJSON(http.StatusOK, dto.CalculationResponseDTO{
				Id: id,
			})
			return
		}
	}

	id, err := h.expressionStorage.Create(expr, calculationRequest.IdempotencyKey)
	if err != nil {
		dto.NewResponseError(http.StatusInternalServerError, err.Error()).Abort(c)
		return
	}

	tokens, err := expression.TokenizeExpression(&expr)
	if err != nil {
		dto.NewResponseError(http.StatusBadRequest, err.Error()).Abort(c)
		return
	}

	root := binary_tree.NewBinaryTree(binary_tree.TokensToNodeArray(tokens))

	err = h.binaryTreeStorage.SaveTree(root, id, -1)
	if err != nil {
		dto.NewResponseError(http.StatusInternalServerError, err.Error()).Abort(c)
		return
	}

	c.IndentedJSON(http.StatusOK, dto.CalculationResponseDTO{
		Id: id,
	})

}

func (h *HTTPHandler) handleExpressionStatusRequest(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		dto.NewResponseError(http.StatusBadRequest, err.Error()).Abort(c)
		return
	}

	if id <= 0 {
		dto.NewResponseError(http.StatusBadRequest, "invalid id").Abort(c)
		return
	}

	statusResponse, err := h.expressionStorage.FindById(id)
	if errors.Is(err, expressions_storage.ErrExpressionNotFound) {
		dto.NewResponseError(http.StatusNotFound, "expression not found").Abort(c)
		return
	}

	if err != nil {
		dto.NewResponseError(http.StatusInternalServerError, err.Error()).Abort(c)
		return
	}

	c.IndentedJSON(http.StatusOK, statusResponse)
}

func (h *HTTPHandler) getAllExpressions(c *gin.Context) {
	expressions, err := h.expressionStorage.FindAll()
	if err != nil {
		dto.NewResponseError(http.StatusInternalServerError, err.Error()).Abort(c)
		return
	}

	c.IndentedJSON(http.StatusOK, expressions)
}

func (h *HTTPHandler) handleWorkerRegister(c *gin.Context) {
	var worker = &dto.WorkerRequestDTO{}

	err := c.Bind(worker)
	if err != nil {
		dto.NewResponseError(http.StatusBadRequest, err.Error()).Abort(c)
		return
	}

	if worker.Url == "" {
		dto.NewResponseError(http.StatusBadRequest, "empty url").Abort(c)
		return
	}

	err = h.workersStorage.Register(worker)
	if err != nil {
		dto.NewResponseError(http.StatusInternalServerError, err.Error()).Abort(c)
		return
	}
}

func (h *HTTPHandler) getAllWorkers(c *gin.Context) {
	workers, err := h.workersStorage.FindAll()
	if err != nil {
		dto.NewResponseError(http.StatusInternalServerError, err.Error()).Abort(c)
		return
	}

	c.IndentedJSON(http.StatusOK, workers)
}
