package main

import (
	"context"
	"fmt"
	"github.com/AleksandrVishniakov/distributed-calculator/daemon/app/internal/servers/grpcsrv"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/AleksandrVishniakov/distributed-calculator/daemon/app/internal/services/executors_pool"
	"github.com/AleksandrVishniakov/distributed-calculator/daemon/app/internal/services/orhestrator_pinger"
)

func main() {
	ctx := context.Background()
	wg := &sync.WaitGroup{}

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

	poolManager := executors_pool.NewManager(executors)
	defer poolManager.Shutdown()

	//handler := handlers.NewHTTPHandler(os.Getenv("ORCHESTRATOR_HOST"), poolManager)
	//server := httpsrv.NewHTTPServer(os.Getenv("HTTP_PORT"), handler.InitRoutes())

	pinger, err := orhestrator_pinger.NewOrchestratorPinger(
		ctx,
		uint64(id),
		os.Getenv("DAEMON_HOST"),
		os.Getenv("ORCHESTRATOR_HOST"),
		executors,
	)

	if err != nil {
		log.Fatal(err)
	}

	pinger.MustPingOrchestrator(
		ctx,
		time.Duration(period)*time.Millisecond,
	)

	gRPCServer := grpc.NewServer()
	grpcsrv.Register(gRPCServer, os.Getenv("ORCHESTRATOR_HOST"), poolManager)

	//go func() {
	//	log.Println("server started on port", os.Getenv("HTTP_PORT"))
	//	if err := server.Run(); err != nil {
	//		server.Shutdown(context.Background())
	//		log.Println(err)
	//	}
	//}()

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
