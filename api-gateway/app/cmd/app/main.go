package main

import (
	"context"
	"fmt"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/dto"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/handlers"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/repositories/expr_tree_repository"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/repositories/expressions_repository"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/repositories/operator_repository"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/repositories/postgres"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/repositories/workers_repository"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/servers"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/servers/grpcsrv"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/services/binary_tree_storage"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/services/expression/expr_tokens"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/services/expressions_storage"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/services/operators_storage"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/services/worker_api"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/services/workers_storage"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/pkg/jwt"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/util/configs"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

func main() {
	wg := &sync.WaitGroup{}

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
	operatorsRepository := operator_repository.NewOperatorsRepository(db)

	workersStorage := workers_storage.NewWorkerStorage(workersRepository)
	expressionStorage := expressions_storage.NewExpressionStorage(expressionsRepository)
	binaryTreeStorage := binary_tree_storage.NewBinaryTreeStorage(binaryTreeRepository)
	operatorsStorage := operators_storage.NewOperatorsStorage(operatorsRepository)

	workerAPI := worker_api.NewGRPCWorkerAPI()

	tokensGenerator := &jwt.TokenGenerator{
		Signature: []byte(os.Getenv("JWT_SIGNATURE")),
		TokenTTL:  12 * time.Hour,
	}

	handler := handlers.NewHTTPHandler(
		expressionStorage,
		binaryTreeStorage,
		workersStorage,
		operatorsStorage,
		workerAPI,
		tokensGenerator,
	)
	server := servers.NewHTTPServer(httpPort, handler.InitRoutes())

	err = operationsInit(operatorsStorage, 500)
	if err != nil {
		log.Fatalf("operators init error: %s", err.Error())
	}

	err = binaryTreeStorage.DeleteAllWorkers()
	if err != nil {
		log.Fatalf("all workers from binary tree deleting error: %s", err.Error())
	}

	_, err = workersStorage.DeleteExpiredWorkers(time.Now())
	if err != nil {
		log.Fatalf("all workers deleting error: %s", err.Error())
	}

	monitorWorkers(workersStorage, binaryTreeStorage)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := server.Run(); err != nil {
			server.Shutdown(context.Background())
			log.Fatalf("server shutted down: %s\n", err.Error())
		}
	}()

	gRPCServer := grpc.NewServer()
	grpcsrv.Register(
		gRPCServer,
		binaryTreeStorage,
		operatorsStorage,
		workersStorage,
		expressionStorage,
		workerAPI,
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		l, err := net.Listen("tcp", fmt.Sprintf(":%s", os.Getenv("GRPC_PORT")))
		if err != nil {
			log.Println(err)
		}
		if err := gRPCServer.Serve(l); err != nil {
			log.Println(err)
		}
	}()

	wg.Wait()
}

func envInit() {
	if err := godotenv.Load(); err != nil {
		log.Print(err)
	}
}

func monitorWorkers(
	workerStorage workers_storage.WorkerStorage,
	binaryTreeStorage binary_tree_storage.BinaryTreeStorage,
) {
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
				workerIds, err := workerStorage.DeleteExpiredWorkers(t.Add(-1 * dur))
				if err != nil {
					log.Fatalf("expired workers deleting error: %s", err.Error())
				}

				if len(workerIds) > 0 {
					err = binaryTreeStorage.DeleteWorkers(workerIds)
					if err != nil {
						log.Fatalf("binary tree cleaning error: %s", err.Error())
					}
				}
			}
		}
	}()
}

func operationsInit(storage operators_storage.OperatorsStorage, defaultDurationMS int) error {
	operations, err := storage.FindAll()
	if err != nil {
		return err
	}

	if len(operations) > 0 {
		return nil
	}

	err = storage.SaveAll([]*dto.OperationDTO{
		{
			OperationType: expr_tokens.Plus,
			DurationMS:    defaultDurationMS,
		},

		{
			OperationType: expr_tokens.Minus,
			DurationMS:    defaultDurationMS,
		},

		{
			OperationType: expr_tokens.Multiply,
			DurationMS:    defaultDurationMS,
		},

		{
			OperationType: expr_tokens.Divide,
			DurationMS:    defaultDurationMS,
		},
	})

	return err
}
