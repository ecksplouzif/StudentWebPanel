package main

import (
	"WebPanel/internal/handlers"
	"WebPanel/internal/middleware"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	mux := http.NewServeMux()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux.HandleFunc("/healthz", handlers.Healthz)
	mux.HandleFunc("/", handlers.NotFound)
	mux.HandleFunc("/login", handlers.Login)
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
