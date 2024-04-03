package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
)

func NewRecoveryMiddleware(log *slog.Logger) func(next http.Handler) http.Handler {
	const src = "http.Recovery"
	log = log.With(
		slog.String("src", src),
	)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					w.WriteHeader(http.StatusInternalServerError)

					log.Error("panic recovered",
						slog.String("error", fmt.Sprintf("%v", err)),
						slog.String("url", r.URL.String()),
					)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
