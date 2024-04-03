package httpsrv

import (
	"context"
	"net"
	"net/http"
	"strconv"
)

type HTTPServer struct {
	httpServer *http.Server
}

func New(ctx context.Context, port int, handler http.Handler) *HTTPServer {
	return &HTTPServer{
		httpServer: &http.Server{
			Addr:    ":" + strconv.Itoa(port),
			Handler: handler,
			BaseContext: func(net.Listener) context.Context {
				return ctx
			},
		},
	}
}

func (s *HTTPServer) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *HTTPServer) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
