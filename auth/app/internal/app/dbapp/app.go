package dbapp

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/AleksandrVishniakov/distributed-calculator/auth/app/pkg/sl"
)

type App struct {
	DB *sql.DB

	log *slog.Logger
	cfg *DBConfigs
}

func New(
	log *slog.Logger,
	cfg *DBConfigs,
) (*App, error) {
	log = log.With(
		slog.Group("cfg",
			slog.String("addr", net.JoinHostPort(cfg.Host, cfg.Port)),
			slog.String("user", cfg.Username),
			slog.String("dbname", cfg.DBName),
		),
	)

	db, err := NewPostgresDB(cfg)
	if err != nil {
		log.Error("failed to start db")
		return nil, err
	}

	db.SetMaxOpenConns(cfg.MaxOpenConnections)
	db.SetMaxIdleConns(cfg.MaxIdleConnections)

	return &App{
		log: log,
		cfg: cfg,
		DB:  db,
	}, nil
}

func (a *App) Healthy() error {
	const src = "dbapp.Healthy"
	log := a.log.With(
		slog.String("src", src),
	)

	log.Debug("ping")

	if err := a.DB.Ping(); err != nil {
		log.Error("ping failed", sl.Err(err))
		return err
	}

	return nil
}

func (a *App) WaitForStart(retries int, interval time.Duration) error {
	const src = "dbapp.WaitForStart"
	log := a.log.With(
		slog.String("src", src),
	)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	var tryCount = 1

	for tryCount <= retries {
		select {
		case <-ticker.C:
			err := a.Healthy()
			if err != nil {
				log.Warn(fmt.Sprintf("try %d/%d failed", tryCount, retries))
				tryCount++
				break
			}

			log.Info("db connected")
			return nil
		}
	}

	log.Error("db is unavailable")
	return errors.New("db is unavailable")
}

func (a *App) Stop() error {
	const src = "dbapp.Stop"
	log := a.log.With(
		slog.String("src", src),
	)

	log.Debug("stopping dbapp...")

	if err := a.DB.Close(); err != nil {
		log.Error("close failed", sl.Err(err))
		return err
	}

	log.Info("dbapp stopped")

	return nil
}
