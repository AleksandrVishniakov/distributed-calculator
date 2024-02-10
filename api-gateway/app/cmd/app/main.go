package main

import (
	"context"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/handlers"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/repositories/expr_tree_repository"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/repositories/expressions_repository"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/repositories/postgres"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/repositories/workers_repository"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/servers"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/services/binary_tree_storage"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/services/expressions_storage"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/services/workers_storage"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/util/configs"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	envInit()

	var httpPort = os.Getenv("HTTP_PORT")
	cfg := configs.MustConfigs()
	cfg.DB.Password = os.Getenv("DB_PASSWORD")

	db, err := postgres.NewPostgresDB(cfg.DB)
	if err != nil {
		log.Fatalf("error while starting postgresql: %s", err.Error())
	}

	expressionsRepository := expressions_repository.NewExpressionsRepository(db)
	binaryTreeRepository := expr_tree_repository.NewExpressionsTreeRepository(db)
	workersRepository := workers_repository.NewWorkersRepository(db)

	workersStorage := workers_storage.NewWorkerStorage(workersRepository)
	expressionStorage := expressions_storage.NewExpressionStorage(expressionsRepository)
	binaryTreeStorage := binary_tree_storage.NewBinaryTreeStorage(binaryTreeRepository)

	handler := handlers.NewHTTPHandler(
		expressionStorage,
		binaryTreeStorage,
		workersStorage,
	)
	server := servers.NewHTTPServer(httpPort, handler.InitRoutes())

	monitorWorkers(workersStorage)

	if err := server.Run(); err != nil {
		server.Shutdown(context.Background())
		log.Fatalf("server shutted down: %s\n", err.Error())
	}
}

func envInit() {
	if err := godotenv.Load(); err != nil {
		log.Print(err)
	}
}

func monitorWorkers(storage workers_storage.WorkerStorage) {
	period, err := strconv.Atoi(os.Getenv("WORKERS_MONITORING_PERIOD_MS"))
	if err != nil {
		log.Fatalf("workers monitoring error: %s", err.Error())
	}

	go func() {
		dur := time.Duration(period) * time.Millisecond
		ticker := time.NewTicker(dur)

		defer ticker.Stop()

		for {
			select {
			case t := <-ticker.C:
				_, err := storage.DeleteExpiredWorkers(t.Add(-1 * dur))
				if err != nil {
					log.Fatalf("workers deleting error: %s", err.Error())
				}
			}
		}
	}()
}
