package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

func NotFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	response := map[string]string{"message": "Not Found"}
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println(err)
	}
}
