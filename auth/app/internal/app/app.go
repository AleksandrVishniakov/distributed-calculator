package app

import (
	"context"
	"log/slog"
	"strconv"
	"sync"
	"time"

	"github.com/AleksandrVishniakov/distributed-calculator/auth/app/internal/app/dbapp"
	"github.com/AleksandrVishniakov/distributed-calculator/auth/app/internal/app/grpcapp"
	"github.com/AleksandrVishniakov/distributed-calculator/auth/app/internal/app/httpapp"
	"github.com/AleksandrVishniakov/distributed-calculator/auth/app/internal/repositories/usersrepo"
	"github.com/AleksandrVishniakov/distributed-calculator/auth/app/internal/servers/httpsrv/handlers"
	"github.com/AleksandrVishniakov/distributed-calculator/auth/app/internal/services/auth"
	"github.com/AleksandrVishniakov/distributed-calculator/auth/app/internal/services/userscache"
	"github.com/AleksandrVishniakov/distributed-calculator/auth/app/pkg/jwt"
	"github.com/AleksandrVishniakov/distributed-calculator/auth/app/utils/configs"

	"google.golang.org/grpc"
)

type App struct {
	ctx context.Context
	log *slog.Logger

	gRPCServer *grpcapp.App
	httpApp    *httpapp.App
	dbApp      *dbapp.App
}

func New(
	ctx context.Context,
	log *slog.Logger,
	cfg *configs.Config,
	getenv func(string) string,
) (*App, error) {
	dbApp, err := dbapp.New(log, &dbapp.DBConfigs{
		Host:               cfg.DB.Host,
		Port:               strconv.Itoa(cfg.DB.Port),
		Username:           cfg.DB.Username,
		DBName:             cfg.DB.DBName,
		Password:           getenv("DB_PASSWORD"),
		SSLMode:            cfg.DB.SSLMode,
		MaxOpenConnections: cfg.DB.MaxConnections,
		MaxIdleConnections: cfg.DB.MaxIdleConnections,
	})

	if err != nil {
		return nil, err
	}

	usersRepository := usersrepo.New(log, dbApp.DB)

	usersCache := userscache.New(log, usersRepository, cfg.Cache.MaxSize, cfg.Cache.TTL)

	tokenGenerator := &jwt.TokenGenerator{
		Signature: []byte(getenv("JWT_SIGNATURE")),
		TokenTTL:  cfg.TokenTTL,
	}

	authService := auth.New(log, usersRepository, usersCache, tokenGenerator)

	gRPCServer := grpc.NewServer()

	grpcApp := grpcapp.New(
		log,
		gRPCServer,
		cfg.GRPC.Port,
		authService,
	)

	httpHandler := handlers.NewHTTPHandler(log, authService)

	httpApp := httpapp.New(ctx, log, cfg.HTTP.Port, httpHandler.Handler())

	return &App{
		ctx:        ctx,
		log:        log,
		gRPCServer: grpcApp,
		httpApp:    httpApp,
		dbApp:      dbApp,
	}, nil
}

func (a *App) Run(wg *sync.WaitGroup) error {
	defer wg.Done()

	const src = "App.Run"
	log := a.log.With(
		slog.String("src", src),
	)

	err := a.dbApp.WaitForStart(3, 3*time.Second)
	if err != nil {
		return nil
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		a.gRPCServer.MustRun()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		a.httpApp.MustRun()
	}()

	log.Info("app started")

	return nil
}

func (a *App) Stop() error {
	const src = "App.Stop"
	log := a.log.With(
		slog.String("src", src),
	)

	log.Debug("stopping app...")

	err := a.dbApp.Stop()
	if err != nil {
		return err
	}

	err = a.httpApp.Stop(a.ctx)
	if err != nil {
		return err
	}

	a.gRPCServer.Stop()

	log.Info("app stopped")

	return nil
}
