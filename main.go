package main

import (
	"fmt"
	"net/http"

	"sync/atomic"

	"time"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	serveMux := http.NewServeMux()

	apiCfg := &apiConfig{}

	serveMux.Handle("/app/", prefixAppStripper(apiCfg.middlewareMetric(http.FileServer(http.Dir(".")))))

	serveMux.Handle("/assets/", http.FileServer(http.Dir("")))

	serveMux.HandleFunc("GET /api/healthz", handlerHealthz)

	serveMux.HandleFunc("GET /admin/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		doc := fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", apiCfg.fileserverHits.Load())
		w.Write([]byte(doc))
	})

	serveMux.HandleFunc("POST /admin/reset", func(w http.ResponseWriter, r *http.Request) {
		apiCfg.fileserverHits.Store(0)
		w.WriteHeader(http.StatusOK)
	})

	server := &http.Server{
		Addr:         ":8080",
		Handler:      serveMux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
