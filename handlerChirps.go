package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/SittingBox/Chirpy/internal/auth"
	"github.com/SittingBox/Chirpy/internal/database"
	"github.com/google/uuid"
)

func handlerChirps(w http.ResponseWriter, r *http.Request, cfg *apiConfig) {
	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not resolve bearer token")
		return
	}

	userID, err := auth.ValidateJWT(bearerToken, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not validate token")
		return
	}

	handlerValidateChirp(w, r, userID, cfg)
}

func handlerValidateChirp(w http.ResponseWriter, r *http.Request, userID uuid.UUID, cfg *apiConfig) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	parms := parameters{}
	err := decoder.Decode(&parms)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Could not decode sent json")
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

	createdChirp, err := cfg.dbQueries.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleanString,
		UserID: userID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error sending chirp")
		return
	}

	type validChirp struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}

	chirp := validChirp{
		ID:        createdChirp.ID,
		CreatedAt: createdChirp.CreatedAt,
		UpdatedAt: createdChirp.UpdatedAt,
		Body:      createdChirp.Body,
		UserID:    createdChirp.UserID,
	}

	respondWithJSON(w, http.StatusCreated, chirp)

}
