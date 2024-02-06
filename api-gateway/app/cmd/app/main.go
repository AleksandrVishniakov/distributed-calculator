package main

import (
	"context"
	"github.com/AleksandrVishniakov/distributed-calculator/api-gateway/app/internal/servers"
	"log"
	"os"
)

func main() {
	var httpPort = os.Getenv("HTTP_PORT")
	
	server := servers.NewHTTPServer(httpPort, nil)
	
	if err := server.Run(); err != nil {
		server.Shutdown(context.Background())
		log.Fatalf("server shutted down: %s\n", err.Error())
	}
}