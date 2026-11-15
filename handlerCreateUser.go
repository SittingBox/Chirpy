package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func handlerCreateUser(w http.ResponseWriter, r *http.Request, cfg *apiConfig) {
	type userEmail struct {
		Email string `json:"email"`
	}
	decoder := json.NewDecoder(r.Body)
	email := userEmail{}
	err := decoder.Decode(&email)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error decoding email")
		return
	}

	type user struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}

	createdUser, err := cfg.dbQueries.CreateUser(r.Context(), email.Email)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error creating user")
		return
	}

	newUser := user{
		ID:        createdUser.ID,
		CreatedAt: createdUser.CreatedAt,
		UpdatedAt: createdUser.UpdatedAt,
		Email:     createdUser.Email,
	}

	fmt.Printf("%+v\n", newUser)

	respondWithJSON(w, 201, newUser)

}
