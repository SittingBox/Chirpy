package main

import "net/http"

func handlerReset(w http.ResponseWriter, r *http.Request, cfg *apiConfig) {
	if cfg.platform != "dev" {
		respondWithError(w, http.StatusForbidden, "Not Authorized")
		return
	}

	err := cfg.dbQueries.ResetUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error reseting users")
		return
	}

	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
}
