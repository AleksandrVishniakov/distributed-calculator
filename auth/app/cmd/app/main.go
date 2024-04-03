package main

import (
	"context"
	stdLog "log"
	"log/slog"
	"os"
	"sync"

	"github.com/AleksandrVishniakov/distributed-calculator/auth/app/internal/app"
	"github.com/AleksandrVishniakov/distributed-calculator/auth/app/utils/configs"
)

func run(
	ctx context.Context,
	log *slog.Logger,
	cfg *configs.Config,
	getenv func(string) string,
) error {
	application, err := app.New(
		ctx,
		log,
		cfg,
		getenv,
	)

	if err != nil {
		return err
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	err = application.Run(wg)
	if err != nil {
		return err
	}

	defer func() {
		err := application.Stop()
		if err != nil {
			stdLog.Fatalf("application stopping error: %s", err.Error())
		}
	}()

	wg.Wait()

	return nil
}

func main() {
	cfg := configs.MustLoad()

	log := logger(cfg.Env)

	ctx, cancel := context.WithCancel(context.Background())

	if err := run(ctx, log, cfg, os.Getenv); err != nil {
		cancel()
		stdLog.Fatalf("server working error: %s", err.Error())
	}
}

func logger(env configs.Environment) *slog.Logger {
	if env != configs.EnvLocal && env != configs.EnvProd {
		stdLog.Printf("unknown environment: %s; set environment to default %s", env, configs.EnvLocal)
		env = configs.EnvLocal
	}

	var log *slog.Logger

	switch env {
	case configs.EnvProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}

	return log
}
