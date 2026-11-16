package main

import (
	"net/http"
	"time"

	"github.com/google/uuid"
)

func handlerRequestChirps(w http.ResponseWriter, r *http.Request, cfg *apiConfig) {
	type Chirp struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}

	chirps, err := cfg.dbQueries.RetrieveChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error fetching chirps")
		return
	}

	var newChirps []Chirp

	for _, chirp := range chirps {
		newChirp := Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		}

		newChirps = append(newChirps, newChirp)
	}

	respondWithJSON(w, http.StatusOK, newChirps)

}
