package handlers

import (
	"log/slog"
	"net/http"
	"time"
)

type writerRecorder struct {
	w      http.ResponseWriter
	body   []byte
	status int
}

func (r *writerRecorder) WriteHeader(status int) {
	r.status = status
	r.w.WriteHeader(status)
}

func (r *writerRecorder) Header() http.Header {
	return r.w.Header()
}

func (r *writerRecorder) Write(bytes []byte) (int, error) {
	r.body = bytes

	return r.w.Write(bytes)
}

func NewLoggerMiddleware(log *slog.Logger) func(next http.Handler) http.Handler {
	const src = "http.Logger"
	log = log.With(
		slog.String("src", src),
	)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var startedAt = time.Now()

			recorder := &writerRecorder{
				w:      w,
				status: http.StatusOK,
				body:   make([]byte, 0),
			}

			next.ServeHTTP(recorder, r)

			var duration = time.Since(startedAt)
			var url = r.URL.String()
			var statusCode = recorder.status
			var method = r.Method

			logFunc := loggingMethod(log, statusCode)

			logFunc("http",
				slog.Group("request",
					slog.String("url", url),
					slog.String("method", method),
				),

				slog.Group("response",
					slog.Int("code", statusCode),
					slog.String("duration", duration.String()),
				),
			)
		})
	}
}

func loggingMethod(log *slog.Logger, status int) func(msg string, args ...any) {
	if status >= 200 && status < 300 {
		return log.Info
	} else {
		return log.Warn
	}
}
