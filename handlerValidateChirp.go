package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.WriteHeader(code)
	encoder := json.NewEncoder(w)
	err := encoder.Encode(&payload)
	if err != nil {
		log.Printf("Error encoding parameters: %s", err)
		return
	}
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type errorReturn struct {
		Error string `json:"error"`
	}

	w.WriteHeader(code)
	parms := errorReturn{
		Error: msg,
	}
	encoder := json.NewEncoder(w)
	err := encoder.Encode(&parms)
	if err != nil {
		log.Printf("Error encoding parameters: %s", err)
		return
	}

}

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	parms := parameters{}
	err := decoder.Decode(&parms)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}

	if len(parms.Body) > 140 {
		respondWithError(w, 400, "Chirp too long")
		return
	} else {
		type validResponse struct {
			Valid bool `json:"valid"`
		}
		resp := validResponse{
			Valid: true,
		}
		respondWithJSON(w, 200, resp)
	}
}
