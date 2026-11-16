package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/SittingBox/Chirpy/internal/database"
	"github.com/google/uuid"
)

func handlerChirps(w http.ResponseWriter, r *http.Request, cfg *apiConfig) {
	handlerValidateChirp(w, r, cfg)
}

func handlerValidateChirp(w http.ResponseWriter, r *http.Request, cfg *apiConfig) {
	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
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
		respondWithError(w, http.StatusBadRequest, "Chirp too long")
		return
	}

	words := strings.Split(parms.Body, " ")

	for i, word := range words {
		lowerWord := strings.ToLower(word)
		if lowerWord == "kerfuffle" || lowerWord == "sharbert" || lowerWord == "fornax" {
			words[i] = "****"
		}
	}

	cleanString := strings.Join(words, " ")

	chirp, err := cfg.dbQueries.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleanString,
		UserID: parms.UserID,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error sending chirp")
		return
	}

	type validChirp struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}

	Chirp := validChirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}

	respondWithJSON(w, http.StatusCreated, Chirp)

}
