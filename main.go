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

	serveMux.Handle("/app/", prefixStripper(apiCfg.middlewareMetric(http.FileServer(http.Dir(".")))))

	serveMux.HandleFunc("/healthz", handlerHealthz)

	serveMux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Hits: %d", apiCfg.fileserverHits.Load())
	})

	serveMux.HandleFunc("/reset", func(w http.ResponseWriter, r *http.Request) {
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
