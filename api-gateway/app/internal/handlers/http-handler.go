package handlers

import (
	"errors"
	"fmt"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/dto"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/handlers/middlewares"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/services/binary_tree"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/services/binary_tree_storage"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/services/expression"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/services/expressions_storage"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/services/operators_storage"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/services/worker_api"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/services/workers_storage"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/pkg/jwt"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/util/calc"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

const tempUserID = 1

type HTTPHandler struct {
	expressionStorage expressions_storage.ExpressionStorage
	binaryTreeStorage binary_tree_storage.BinaryTreeStorage
	workersStorage    workers_storage.WorkerStorage
	operatorsStorage  operators_storage.OperatorsStorage
	workerAPI         worker_api.WorkerAPI

	tokensGenerator *jwt.TokenGenerator
}

func NewHTTPHandler(
	expressionStorage expressions_storage.ExpressionStorage,
	binaryTreeStorage binary_tree_storage.BinaryTreeStorage,
	workersStorage workers_storage.WorkerStorage,
	operatorsStorage operators_storage.OperatorsStorage,
	workerAPI worker_api.WorkerAPI,
	tokensGenerator *jwt.TokenGenerator,
) *HTTPHandler {
	return &HTTPHandler{
		expressionStorage: expressionStorage,
		binaryTreeStorage: binaryTreeStorage,
		workersStorage:    workersStorage,
		operatorsStorage:  operatorsStorage,
		workerAPI:         workerAPI,
		tokensGenerator:   tokensGenerator,
	}
}

func (h *HTTPHandler) InitRoutes() *gin.Engine {
	router := gin.Default()

	router.Use(middlewares.CORSHeaders())

	jwtAuth := middlewares.NewJWTAuthMiddleware(h.tokensGenerator)

	api := router.Group("/api")
	{
		api.GET("/operators", h.getAllOperations)
		api.POST("/operators", h.saveAllOperations)

		api.POST("/expression", jwtAuth(), h.calculateExpression)
		api.GET("/expressions", jwtAuth(), h.getAllExpressions)
		api.GET("/expression/:id", h.handleExpressionStatusRequest)

		api.POST("/task/:id/result", h.handleTaskResult)
		api.POST("/task/:id/status", h.handleTaskStarting)

		api.POST("/worker", h.handleWorkerRegister)
		api.GET("/worker/:id/tasks", h.handleWorkerTasks)
		api.GET("/workers", h.getAllWorkers)
	}

	return router
}

func (h *HTTPHandler) calculateExpression(c *gin.Context) {
	userID, err := userID(c)
	if err != nil {
		dto.NewResponseError(http.StatusBadRequest, err.Error()).Abort(c)
		return
	}

	var calculationRequest dto.CalculationRequestDTO

	err = c.BindJSON(&calculationRequest)
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

	expressionId, err := h.expressionStorage.Create(expr, userID, calculationRequest.IdempotencyKey)
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

	taskId, err := h.binaryTreeStorage.SaveTree(root, userID, expressionId, -1, true)
	if err != nil {
		dto.NewResponseError(http.StatusInternalServerError, err.Error()).Abort(c)
		return
	}

	err = calc.StartCalculating(
		c,
		taskId,
		h.binaryTreeStorage,
		h.operatorsStorage,
		h.workersStorage,
		h.expressionStorage,
		h.workerAPI,
	)
	if err != nil {
		dto.NewResponseError(http.StatusInternalServerError, err.Error()).Abort(c)
		return
	}

	c.IndentedJSON(http.StatusOK, dto.CalculationResponseDTO{
		Id: expressionId,
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
	userID, err := userID(c)
	if err != nil {
		dto.NewResponseError(http.StatusBadRequest, err.Error()).Abort(c)
		return
	}

	expressions, err := h.expressionStorage.FindAllByUserID(userID)
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

	exists, err := h.workersStorage.Register(worker)
	if err != nil {
		dto.NewResponseError(http.StatusInternalServerError, err.Error()).Abort(c)
		return
	}

	if exists {
		return
	}

	err = calc.CalculateAll(
		c,
		h.binaryTreeStorage,
		h.operatorsStorage,
		h.workersStorage,
		h.expressionStorage,
		h.workerAPI,
	)
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

func (h *HTTPHandler) handleTaskStarting(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		dto.NewResponseError(http.StatusBadRequest, err.Error()).Abort(c)
		return
	}

	err = h.binaryTreeStorage.MarkAsCalculating(id)
	if err != nil {
		dto.NewResponseError(http.StatusInternalServerError, err.Error()).Abort(c)
		return
	}

	node, err := h.binaryTreeStorage.FindById(id)
	if err != nil {
		dto.NewResponseError(http.StatusInternalServerError, err.Error()).Abort(c)
		return
	}

	err = h.expressionStorage.MarkAsCalculating(node.ExpressionId)
	if err != nil {
		dto.NewResponseError(http.StatusInternalServerError, err.Error()).Abort(c)
		return
	}
}

func (h *HTTPHandler) handleTaskResult(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		dto.NewResponseError(http.StatusBadRequest, err.Error()).Abort(c)
		return
	}

	var calculationResult = &dto.CalculationResultDTO{}

	err = c.BindJSON(calculationResult)
	if err != nil {
		dto.NewResponseError(http.StatusBadRequest, err.Error()).Abort(c)
		return
	}

	err = h.binaryTreeStorage.SaveResult(id, calculationResult.Result)
	if err != nil {
		dto.NewResponseError(http.StatusInternalServerError, err.Error()).Abort(c)
		return
	}

	node, err := h.binaryTreeStorage.FindById(id)
	if err != nil {
		dto.NewResponseError(http.StatusInternalServerError, err.Error()).Abort(c)
		return
	}

	if node.ParentId == -1 {
		err = h.expressionStorage.SaveResult(node.ExpressionId, calculationResult.Result)
		if err != nil {
			dto.NewResponseError(http.StatusInternalServerError, err.Error()).Abort(c)
		}

		return
	}

	//err = h.startCalculating(node.ParentId)
	//if err != nil {
	//	dto.NewResponseError(http.StatusInternalServerError, err.Error()).Abort(c)
	//	return
	//}

	err = calc.CalculateAll(
		c,
		h.binaryTreeStorage,
		h.operatorsStorage,
		h.workersStorage,
		h.expressionStorage,
		h.workerAPI,
	)
	if err != nil {
		dto.NewResponseError(http.StatusInternalServerError, err.Error()).Abort(c)
	}
}

func (h *HTTPHandler) getAllOperations(c *gin.Context) {
	operations, err := h.operatorsStorage.FindAll()
	if err != nil {
		dto.NewResponseError(http.StatusInternalServerError, err.Error()).Abort(c)
		return
	}

	c.IndentedJSON(http.StatusOK, operations)
}

func (h *HTTPHandler) saveAllOperations(c *gin.Context) {
	var operations []*dto.OperationDTO

	err := c.BindJSON(&operations)
	if err != nil {
		dto.NewResponseError(http.StatusBadRequest, err.Error()).Abort(c)
		return
	}

	for _, operation := range operations {
		if operation.DurationMS < 0 {
			dto.NewResponseError(http.StatusBadRequest, "invalid operation duration time").Abort(c)
			return
		}
	}

	err = h.operatorsStorage.SaveAll(operations)
	if err != nil {
		dto.NewResponseError(http.StatusInternalServerError, err.Error()).Abort(c)
		return
	}
}

func (h *HTTPHandler) handleWorkerTasks(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		dto.NewResponseError(http.StatusBadRequest, err.Error()).Abort(c)
		return
	}

	tasks, err := h.binaryTreeStorage.FindByWorkerId(id)
	if err != nil {
		dto.NewResponseError(http.StatusInternalServerError, err.Error()).Abort(c)
		return
	}

	c.IndentedJSON(http.StatusOK, tasks)
}

func userID(c *gin.Context) (uint64, error) {
	userID, err := strconv.ParseUint(fmt.Sprintf("%v", c.Value("user_id")), 10, 64)
	if err != nil {
		return 0, err
	}

	return userID, nil

	//return tempUserID, nil
}
