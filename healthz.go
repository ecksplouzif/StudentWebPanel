package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

func healthz(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{"message": "OK"}
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println(err)
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.HandleFunc("/healthz", healthz)
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
