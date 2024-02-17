package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

var httpServer = &http.Server{
	Addr:              ":" + os.Getenv("HTTP_PORT"),
	ReadTimeout:       10 * time.Second,
	ReadHeaderTimeout: 10 * time.Second,
	WriteTimeout:      10 * time.Second,
	IdleTimeout:       10 * time.Second,
}

const (
	indexFilePath  = "web/app/build/index.html"
	staticFilesDir = "web/app/build"
)

func main() {
	mux := http.NewServeMux()

	mux.Handle("/static/", http.FileServer(http.Dir(staticFilesDir)))
	mux.HandleFunc("/", parsePage)

	httpServer.Handler = mux

	log.Println("server started on port", os.Getenv("HTTP_PORT"))
	if err := httpServer.ListenAndServe(); err != nil {
		log.Fatalf("server working error: %s", err.Error())
	}
}

func parsePage(w http.ResponseWriter, r *http.Request) {
	var pageHost = struct {
		Host string
	}{
		Host: os.Getenv("ORCHESTRATOR_HOST"),
	}

	html, err := template.ParseFiles(indexFilePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = html.Execute(w, pageHost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
