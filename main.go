package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/SittingBox/Chirpy/internal/database"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
	platform       string
}

func main() {
	godotenv.Load()

	dbURL := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbURL)

	dbQueries := database.New(db)

	platform := os.Getenv("PLATFORM")
	fmt.Println("PLATFORM loaded as:", platform)

	apiCfg := &apiConfig{
		dbQueries: dbQueries,
		platform:  platform,
	}

	serveMux := http.NewServeMux()

	serveMux.Handle("/app/", prefixAppStripper(apiCfg.middlewareMetric(http.FileServer(http.Dir(".")))))

	serveMux.Handle("/assets/", http.FileServer(http.Dir("")))

	serveMux.HandleFunc("GET /api/healthz", handlerHealthz)

	serveMux.HandleFunc("GET /admin/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		doc := fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", apiCfg.fileserverHits.Load())
		w.Write([]byte(doc))
	})

	serveMux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)

	serveMux.HandleFunc("POST /admin/reset", func(w http.ResponseWriter, r *http.Request) {
		if apiCfg.platform != "dev" {
			respondWithError(w, http.StatusForbidden, "Not Authorized")
			return
		}

		err = apiCfg.dbQueries.ResetUsers(r.Context())
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Error reseting users")
			return
		}

		apiCfg.fileserverHits.Store(0)
		w.WriteHeader(http.StatusOK)
	})

	serveMux.HandleFunc("POST /api/users", func(w http.ResponseWriter, r *http.Request) {
		type userEmail struct {
			Email string `json:"email"`
		}

		decoder := json.NewDecoder(r.Body)
		email := userEmail{}
		err = decoder.Decode(&email)
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

		createdUser, err := apiCfg.dbQueries.CreateUser(r.Context(), email.Email)
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
