package handlers

import (
	"encoding/json"
	"github.com/AleksandrVishniakov/distributed-calculator/daemon/app/internal/dto"
	"github.com/AleksandrVishniakov/distributed-calculator/daemon/app/internal/services/executors_pool"
	"net/http"
)

type HTTPHandler struct {
	orchestratorHost string
	executorsPool    *executors_pool.ExecutorsPool
}

func NewHTTPHandler(
	orchestratorHost string,
	executorsPool *executors_pool.ExecutorsPool,
) *HTTPHandler {
	return &HTTPHandler{
		orchestratorHost: orchestratorHost,
		executorsPool:    executorsPool,
	}
}

func (h *HTTPHandler) InitRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/task", h.calculate)

	return mux
}

func (h *HTTPHandler) calculate(w http.ResponseWriter, r *http.Request) {
	var requestDTO = &dto.CalculationRequestDTO{}

	err := json.NewDecoder(r.Body).Decode(&requestDTO)
	if err != nil {
		dto.NewResponseError(http.StatusBadRequest, err.Error()).Abort(w)
		return
	}

	executor := executors_pool.NewCalculationExecutor(requestDTO, h.orchestratorHost)

	h.executorsPool.Run(executor)
}
