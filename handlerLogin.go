package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/google/uuid"
)

func handlerLogin(w http.ResponseWriter, r *http.Request, cfg *apiConfig) {
	type login struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	loginWithEmail := login{}
	err := decoder.Decode(&loginWithEmail)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error decoding email")
		return
	}

	loginUser, err := cfg.dbQueries.Login(r.Context(), loginWithEmail.Email)

	type user struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}

	userLoggedin := user{
		ID:        loginUser.ID,
		CreatedAt: loginUser.CreatedAt,
		UpdatedAt: loginUser.UpdatedAt,
		Email:     loginUser.Email,
	}

	correctPass, err := argon2id.ComparePasswordAndHash(loginWithEmail.Password, loginUser.HashedPasswords)

	if !correctPass {
		respondWithError(w, http.StatusUnauthorized, "Wrong email or password")
		return
	}

	respondWithJSON(w, http.StatusOK, userLoggedin)
}
