package rest

import (
	"encoding/json"
	"log"
	"net/http"
)

func RespondWithError(w http.ResponseWriter, code int, msg string, err error) {
	if err != nil {
		log.Println(err)
	}

	if code > 499 {
		log.Printf("Responding with error: %s", msg)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}

	RespondWithJson(w, code, errorResponse{Error: msg})

}

func RespondWithJson(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-type", "application/json")
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshaling JSON: %s", err)
	}

	w.WriteHeader(code)
	w.Write(data)

}
