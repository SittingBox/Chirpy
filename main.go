package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/SittingBox/Chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
	platform       string
	secret         string
}

func main() {
	godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}

	db, err := sql.Open("postgres", dbURL)

	dbQueries := database.New(db)

	secret := os.Getenv("SECRET")
	if secret == "" {
		log.Fatal("Secret must be set")
	}

	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("Platform must be set")
	}

	apiCfg := &apiConfig{
		dbQueries: dbQueries,
		platform:  platform,
		secret:    secret,
	}

	serveMux := http.NewServeMux()

	fmt.Println("Server is starting!")

	serveMux.Handle("/app/", prefixAppStripper(apiCfg.middlewareMetric(http.FileServer(http.Dir(".")))))

	serveMux.Handle("/assets/", http.FileServer(http.Dir("")))

	serveMux.HandleFunc("GET /api/healthz", handlerHealthz)

	serveMux.HandleFunc("GET /admin/metrics", func(w http.ResponseWriter, r *http.Request) {
		handlerMetrics(w, apiCfg)
	})

	serveMux.HandleFunc("POST /admin/reset", func(w http.ResponseWriter, r *http.Request) {
		handlerReset(w, r, apiCfg)
	})

	serveMux.HandleFunc("POST /api/users", func(w http.ResponseWriter, r *http.Request) {
		handlerCreateUser(w, r, apiCfg)
	})

	serveMux.HandleFunc("POST /api/login", func(w http.ResponseWriter, r *http.Request) {
		handlerLogin(w, r, apiCfg)
	})

	serveMux.HandleFunc("POST /api/chirps", func(w http.ResponseWriter, r *http.Request) {
		handlerChirps(w, r, apiCfg)
	})

	serveMux.HandleFunc("GET /api/chirps", func(w http.ResponseWriter, r *http.Request) {
		handlerRequestChirps(w, r, apiCfg)
	})

	serveMux.HandleFunc("GET /api/chirps/{chirpID}", func(w http.ResponseWriter, r *http.Request) {
		handlerSingleChirpResponse(w, r, apiCfg)
	})

	server := &http.Server{
		Addr:         ":8080",
		Handler:      serveMux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
