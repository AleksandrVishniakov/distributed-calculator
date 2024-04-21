package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/AleksandrVishniakov/distributed-calculator/daemon/app/internal/dto"
	"github.com/AleksandrVishniakov/distributed-calculator/daemon/app/internal/services/executors_pool"
)

type HTTPHandler struct {
	orchestratorHost string
	poolManager      *executors_pool.PoolManager
}

func NewHTTPHandler(
	orchestratorHost string,
	poolManager *executors_pool.PoolManager,
) *HTTPHandler {
	return &HTTPHandler{
		orchestratorHost: orchestratorHost,
		poolManager:      poolManager,
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

	pool := h.poolManager.Pool(requestDTO.UserID)

	executor, err := executors_pool.NewCalculationExecutor(r.Context(), requestDTO, h.orchestratorHost)
	if err != nil {
		dto.NewResponseError(http.StatusInternalServerError, err.Error()).Abort(w)
		return
	}

	pool.Run(executor)
}
