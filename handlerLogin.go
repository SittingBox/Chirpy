package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/SittingBox/Chirpy/internal/auth"
	"github.com/alexedwards/argon2id"
	"github.com/google/uuid"
)

func handlerLogin(w http.ResponseWriter, r *http.Request, cfg *apiConfig) {
	type login struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}

	decoder := json.NewDecoder(r.Body)
	params := login{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error decoding parameters")
		return
	}

	if params.ExpiresInSeconds > 3600 || params.ExpiresInSeconds <= 0 {
		params.ExpiresInSeconds = 3600
	}

	loginUser, err := cfg.dbQueries.Login(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Wrong email or password")
		return
	}

	correctPass, err := argon2id.ComparePasswordAndHash(params.Password, loginUser.HashedPasswords)
	if !correctPass || err != nil {
		respondWithError(w, http.StatusUnauthorized, "Wrong email or password")
		return
	}

	expiresIn := time.Second * time.Duration(params.ExpiresInSeconds)

	token, err := auth.MakeJWT(loginUser.ID, cfg.secret, expiresIn)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error setting up token")
		return
	}

	type response struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
		Token     string    `json:"token"`
	}

	userLoggedin := response{
		ID:        loginUser.ID,
		CreatedAt: loginUser.CreatedAt,
		UpdatedAt: loginUser.UpdatedAt,
		Email:     loginUser.Email,
		Token:     token,
	}

	respondWithJSON(w, http.StatusOK, userLoggedin)
}
