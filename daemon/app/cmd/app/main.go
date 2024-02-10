package main

import (
	"github.com/AleksandrVishniakov/distributed-calculator/daemon/app/internal/handlers"
	"github.com/AleksandrVishniakov/distributed-calculator/daemon/app/internal/servers"
	"log"
	"os"
)

func main() {
	handler := handlers.NewHTTPHandler()
	server := servers.NewHTTPServer(os.Getenv("HTTP_PORT"), handler.InitRoutes())

	if err := server.Run(); err != nil {
		server.Shutdown()
		log.Fatal(err)
	}
}
