package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/dto"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/handlers/middlewares"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/services/binary_tree"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/services/binary_tree_storage"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/services/expression"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/services/expression/expr_tokens"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/services/expressions_storage"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/services/operators_storage"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/services/statuses"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/services/worker_api"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/services/workers_storage"

	"github.com/gin-gonic/gin"
)

type HTTPHandler struct {
	expressionStorage expressions_storage.ExpressionStorage
	binaryTreeStorage binary_tree_storage.BinaryTreeStorage
	workersStorage    workers_storage.WorkerStorage
	operatorsStorage  operators_storage.OperatorsStorage
	workerAPI         worker_api.WorkerAPI
}

func NewHTTPHandler(
	expressionStorage expressions_storage.ExpressionStorage,
	binaryTreeStorage binary_tree_storage.BinaryTreeStorage,
	workersStorage workers_storage.WorkerStorage,
	operatorsStorage operators_storage.OperatorsStorage,
	workerAPI worker_api.WorkerAPI,
) *HTTPHandler {
	return &HTTPHandler{
		expressionStorage: expressionStorage,
		binaryTreeStorage: binaryTreeStorage,
		workersStorage:    workersStorage,
		operatorsStorage:  operatorsStorage,
		workerAPI:         workerAPI,
	}
}

func (h *HTTPHandler) InitRoutes() *gin.Engine {
	router := gin.Default()

	router.Use(middlewares.CORSHeaders())

	api := router.Group("/api")
	{
		api.GET("/operators", h.getAllOperations)
		api.POST("/operators", h.saveAllOperations)

		api.POST("/expression", h.calculateExpression)
		api.GET("/expressions", h.getAllExpressions)
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

	expressionId, err := h.expressionStorage.Create(expr, calculationRequest.IdempotencyKey)
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

	taskId, err := h.binaryTreeStorage.SaveTree(root, expressionId, -1, true)
	if err != nil {
		dto.NewResponseError(http.StatusInternalServerError, err.Error()).Abort(c)
		return
	}

	err = h.startCalculating(taskId)
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

	exists, err := h.workersStorage.Register(worker)
	if err != nil {
		dto.NewResponseError(http.StatusInternalServerError, err.Error()).Abort(c)
		return
	}

	if exists {
		return
	}

	err = h.calculateAll()
	if err != nil {
		dto.NewResponseError(http.StatusInternalServerError, err.Error()).Abort(c)
		return
	}
}

func (h *HTTPHandler) calculateAll() error {
	ids, err := h.binaryTreeStorage.FindUncalculated()
	if err != nil {
		return err
	}

	for _, id := range ids {
		err = h.startCalculating(id)
		if err != nil {
			return err
		}
	}

	return nil
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

	err = h.calculateAll()
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

func (h *HTTPHandler) startCalculating(id int) error {
	node, err := h.binaryTreeStorage.FindById(id)
	if err != nil {
		return err
	}

	if node.Status != statuses.Created {
		return nil
	}

	nodes, err := h.binaryTreeStorage.FindByParentId(id)
	if err != nil {
		return err
	}

	if len(nodes) == 0 {
		return nil
	}

	left := nodes[0]
	right := nodes[1]

	if left.Status == statuses.Finished && right.Status == statuses.Finished {
		operations, err := h.operatorsStorage.FindAll()
		if err != nil {
			return err
		}

		operationDuration, err := getOperationDuration(operations, expr_tokens.OperationType(node.OperationType))
		if err != nil {
			return err
		}

		worker, err := h.workersStorage.FindFreeWorker()
		if err != nil {
			return err
		}

		if worker == nil {
			return nil
		}

		var operation = expr_tokens.OperationType(node.OperationType)

		if operation == expr_tokens.Divide && right.Result == 0 {
			err = h.expressionStorage.MarkAsFailed(node.ExpressionId)
			if err != nil {
				return err
			}

			err = h.binaryTreeStorage.MarkAsFailed(node.Id)
			if err != nil {
				return err
			}

			return nil
		}

		err = h.workerAPI.Calculate(worker.Url, &worker_api.CalculationRequestDTO{
			Id:        id,
			First:     left.Result,
			Second:    right.Result,
			Operation: operation,
			Duration:  operationDuration,
		})
		if err != nil {
			return err
		}

		err = h.binaryTreeStorage.SaveWorker(id, worker.Id)
		if err != nil {
			return err
		}
	} else {
		if left.Status == statuses.Created {
			err := h.startCalculating(left.Id)
			if err != nil {
				return err
			}
		}

		if right.Status == statuses.Created {
			err := h.startCalculating(right.Id)
			if err != nil {
				return err
			}
		}
	}

	return nil
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

func getOperationDuration(operations []*dto.OperationDTO, operationType expr_tokens.OperationType) (time.Duration, error) {
	for _, operation := range operations {
		if operation.OperationType == operationType {
			return time.Duration(operation.DurationMS) * time.Millisecond, nil
		}
	}

	return 0, errors.New("operation not found")
}
