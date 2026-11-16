package main

import (
	"net/http"
	"time"

	"github.com/google/uuid"
)

func handlerSingleChirpResponse(w http.ResponseWriter, r *http.Request, cfg *apiConfig) {
	chirpID := r.PathValue("chirpID")

	chirpUUID, err := uuid.Parse(chirpID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error getting ID for chirp")
		return
	}

	foundChirp, err := cfg.dbQueries.FetchSingleChirp(r.Context(), chirpUUID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Error finding chirp")
		return
	}

	type Chirp struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}

	chirpJSON := Chirp{
		ID:        foundChirp.ID,
		CreatedAt: foundChirp.CreatedAt,
		UpdatedAt: foundChirp.UpdatedAt,
		Body:      foundChirp.Body,
		UserID:    foundChirp.UserID,
	}

	respondWithJSON(w, http.StatusOK, chirpJSON)
}
