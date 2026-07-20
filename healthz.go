package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
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
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
