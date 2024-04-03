package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/AleksandrVishniakov/distributed-calculator/auth/app/internal/servers"
	"github.com/AleksandrVishniakov/distributed-calculator/auth/app/internal/servers/httpsrv"
	"github.com/AleksandrVishniakov/distributed-calculator/auth/app/internal/services/auth"
	"github.com/AleksandrVishniakov/distributed-calculator/auth/app/pkg/parser"
	"github.com/AleksandrVishniakov/distributed-calculator/auth/app/pkg/sl"
)

type Auth interface {
	RegisterNewUser(
		ctx context.Context,
		login string,
		password string,
	) (userID uint64, err error)

	Login(
		ctx context.Context,
		login string,
		password string,
	) (token string, err error)
}

type HTTPHandler struct {
	log  *slog.Logger
	auth Auth
}

func NewHTTPHandler(
	log *slog.Logger,
	auth Auth,
) *HTTPHandler {
	return &HTTPHandler{
		log:  log,
		auth: auth,
	}
}

func (h *HTTPHandler) Handler() http.Handler {
	mux := http.NewServeMux()
	logger := NewLoggerMiddleware(h.log)
	recovery := NewRecoveryMiddleware(h.log)

	mux.Handle("POST /api/v1/register", Errors(h.Register))
	mux.Handle("POST /api/v1/login", Errors(h.Login))

	return recovery(logger(CORS(mux)))
}

func (h *HTTPHandler) Register(w http.ResponseWriter, r *http.Request) (statusCode int, err error) {
	const src = "HTTPHandler.Register"
	log := h.log.With(
		"src", src,
	)

	user, err := parser.DecodeValid[*httpsrv.UserRequestDTO](r.Body)
	if err != nil {
		return http.StatusBadRequest, err
	}

	id, err := h.auth.RegisterNewUser(r.Context(), user.Login, user.Password)
	if err != nil {
		if errors.Is(err, auth.ErrUserAlreadyExists) {
			return http.StatusConflict, errors.New(servers.MsgUserAlreadyExists)
		}

		return http.StatusInternalServerError, errors.New(servers.MsgInternalError)
	}

	err = parser.EncodeResponse(w, &httpsrv.RegisterResponseDTO{
		ID: id,
	}, http.StatusOK)

	if err != nil {
		log.Error("failed to encode response", sl.Err(err))
		return http.StatusInternalServerError, errors.New(servers.MsgInternalError)
	}

	return http.StatusOK, nil
}

func (h *HTTPHandler) Login(w http.ResponseWriter, r *http.Request) (statusCode int, err error) {
	const src = "HTTPHandler.Register"
	log := h.log.With(
		"src", src,
	)

	user, err := parser.DecodeValid[*httpsrv.UserRequestDTO](r.Body)
	if err != nil {
		return http.StatusBadRequest, err
	}

	token, err := h.auth.Login(r.Context(), user.Login, user.Password)
	if err != nil {
		if errors.Is(err, auth.ErrUserNotFound) {
			return http.StatusNotFound, errors.New(servers.MsgUserNotFound)
		}
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return http.StatusBadRequest, errors.New(servers.MsgInvalidCredentials)
		}

		return http.StatusInternalServerError, errors.New(servers.MsgInternalError)
	}

	err = parser.EncodeResponse(w, &httpsrv.LoginResponseDTO{
		Token: token,
	}, http.StatusOK)

	if err != nil {
		log.Error("failed to encode response", sl.Err(err))
		return http.StatusInternalServerError, errors.New(servers.MsgInternalError)
	}

	return http.StatusOK, nil
}
