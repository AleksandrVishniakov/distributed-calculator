package httpapp

import (
	"context"
	"github.com/AleksandrVishniakov/distributed-calculator/auth/app/pkg/sl"
	"log/slog"
	"net"
	"net/http"
	"strconv"

	"github.com/AleksandrVishniakov/distributed-calculator/auth/app/internal/servers/httpsrv"
)

type App struct {
	log *slog.Logger

	server *httpsrv.HTTPServer
}

func New(
	ctx context.Context,
	log *slog.Logger,

	port int,
	handler http.Handler,
) *App {
	log = log.With(
		slog.String("addr", net.JoinHostPort("localhost", strconv.Itoa(port))),
	)

	return &App{
		log:    log,
		server: httpsrv.New(ctx, port, handler),
	}
}

func (a *App) MustRun() {
	const src = "httpapp.MustRun"
	log := a.log.With(
		slog.String("src", src),
	)

	if err := a.Run(); err != nil {
		log.Error("server failed", sl.Err(err))
		panic(err)
	}
}

func (a *App) Run() error {
	const src = "httpapp.Run"
	log := a.log.With(
		slog.String("src", src),
	)

	log.Info("starting http server...")

	return a.server.Run()
}

func (a *App) Stop(ctx context.Context) error {
	const src = "httpapp.Stop"
	log := a.log.With(
		slog.String("src", src),
	)

	log.Debug("stopping http server...")

	err := a.server.Shutdown(ctx)
	if err != nil {
		log.Error("failed to stop http server", sl.Err(err))
		return err
	}

	log.Info("http server stopped")
	return nil
}
