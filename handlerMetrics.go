package main

import (
	"fmt"
	"net/http"
)

func handlerMetrics(w http.ResponseWriter, cfg *apiConfig) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	doc := fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", cfg.fileserverHits.Load())
	w.Write([]byte(doc))
}
