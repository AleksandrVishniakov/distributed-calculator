package main

import (
	"context"
	"github.com/AleksandrVishniakov/distributed-calculator/daemon/app/internal/handlers"
	"github.com/AleksandrVishniakov/distributed-calculator/daemon/app/internal/servers"
	"github.com/AleksandrVishniakov/distributed-calculator/daemon/app/internal/services/executors_pool"
	"github.com/AleksandrVishniakov/distributed-calculator/daemon/app/internal/services/orhestrator_pinger"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	id, err := strconv.Atoi(os.Getenv("DAEMON_ID"))
	if err != nil {
		log.Fatal(err)
	}

	executors, err := strconv.Atoi(os.Getenv("MAX_GOROUTINES"))
	if err != nil {
		log.Fatal(err)
	}

	period, err := strconv.Atoi(os.Getenv("PING_PERIOD_MS"))
	if err != nil {
		log.Fatal(err)
	}

	pool := executors_pool.NewExecutorsPool(executors)
	defer pool.Shutdown()

	handler := handlers.NewHTTPHandler(os.Getenv("ORCHESTRATOR_HOST"), pool)
	server := servers.NewHTTPServer(os.Getenv("HTTP_PORT"), handler.InitRoutes())

	pinger, err := orhestrator_pinger.NewOrchestratorPinger(
		id,
		os.Getenv("DAEMON_HOST"),
		executors,
		os.Getenv("ORCHESTRATOR_HOST"),
	)

	if err != nil {
		log.Fatal(err)
	}

	pinger.MustPingOrchestrator(
		time.Duration(period) * time.Millisecond,
	)

	log.Println("server started on port", os.Getenv("HTTP_PORT"))
	if err := server.Run(); err != nil {
		server.Shutdown(context.Background())
		log.Fatal(err)
	}
}
