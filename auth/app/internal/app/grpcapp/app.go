package grpcapp

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/AleksandrVishniakov/distributed-calculator/auth/app/internal/servers/grpcsrv"
	"github.com/AleksandrVishniakov/distributed-calculator/auth/app/pkg/e"

	"google.golang.org/grpc"
)

type App struct {
	log    *slog.Logger
	server *grpc.Server
	port   int
}

func New(
	log *slog.Logger,
	server *grpc.Server,
	port int,
	authService grpcsrv.Auth,
) *App {
	grpcsrv.Register(server, authService)

	return &App{
		log:    log,
		server: server,
		port:   port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const src = "grpcapp.Run"

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return e.WrapErr(err, src)
	}

	a.log.Info("grpc server started", slog.String("addr", l.Addr().String()))

	if err := a.server.Serve(l); err != nil {
		return e.WrapErr(err, src)
	}

	return nil
}

func (a *App) Stop() {
	const src = "grpcapp.Stop"

	a.log.With(slog.String("op", src)).
		Info("stopping gRPC server", slog.Int("port", a.port))

	a.server.GracefulStop()
}
