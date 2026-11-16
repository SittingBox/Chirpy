package main

import (
	"encoding/json"
	"log"
	"net/http"
)

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
