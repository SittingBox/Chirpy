package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/SittingBox/Chirpy/internal/auth"
	"github.com/SittingBox/Chirpy/internal/database"
	"github.com/google/uuid"
)

func handlerCreateUser(w http.ResponseWriter, r *http.Request, cfg *apiConfig) {
	type login struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	Login := login{}
	err := decoder.Decode(&Login)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error decoding email")
		return
	}

	hashedPass, err := auth.HashPassword(Login.Password)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error hashing password")
	}

	createdUser, err := cfg.dbQueries.CreateUser(r.Context(), database.CreateUserParams{
		Email:           Login.Email,
		HashedPasswords: hashedPass,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error creating user")
		return
	}

	type user struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}

	newUser := user{
		ID:        createdUser.ID,
		CreatedAt: createdUser.CreatedAt,
		UpdatedAt: createdUser.UpdatedAt,
		Email:     createdUser.Email,
	}

	respondWithJSON(w, 201, newUser)

}
