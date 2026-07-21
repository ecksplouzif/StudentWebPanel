package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

func Healthz(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{"message": "OK"}
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println(err)
	}
}
