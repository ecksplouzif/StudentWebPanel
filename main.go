package main

import (
	"WebPanel/internal/handlers"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.HandleFunc("/healthz", handlers.Healthz)
	serv := &http.Server{
		Addr:           ":" + port,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	err := serv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
