package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"WebPanel/internal/handlers"
	"WebPanel/internal/middleware"
)

func main() {

	mux := http.NewServeMux()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux.HandleFunc("/healthz", handlers.Healthz)
	mux.HandleFunc("/", handlers.NotFound)
	serv := &http.Server{
		Addr:         ":" + port,
		Handler:      middleware.Logging(mux),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	err := serv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
