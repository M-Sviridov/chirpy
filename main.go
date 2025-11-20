package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/M-Sviridov/chirpy/internal/database"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
}

func main() {
	const listenAddr = ":8080"
	const filePathRoot = "."

	godotenv.Load()
	dbUrl := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		os.Exit(1)
	}
	dbQueries := database.New(db)

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
		platform:       platform,
	}

	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(filePathRoot)))))
	mux.HandleFunc("GET /admin/metrics", apiCfg.handleNumRequests)
	mux.HandleFunc("POST /admin/reset", apiCfg.handleReset)
	mux.HandleFunc("POST /api/validate_chirp", apiCfg.handleValidateChirp)
	mux.HandleFunc("GET /api/healthz", handleReadiness)
	mux.HandleFunc("POST /api/users", apiCfg.handleCreateUser)

	srv := http.Server{
		Handler: mux,
		Addr:    listenAddr,
	}

	srv.ListenAndServe()
}

func handleReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func (cfg *apiConfig) handleNumRequests(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	body := fmt.Sprintf(`
	<html>
		<body>
			<h1>Welcome, Chirpy Admin</h1>
			<p>Chirpy has been visited %d times!</p>
		</body>
	</html>`, cfg.fileserverHits.Load())
	w.Write([]byte(body))
}

func (cfg *apiConfig) handleReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, http.StatusForbidden, "cannot run reset on non dev platform")
	} else {
		cfg.fileserverHits.Store(0)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Reset total hits to 0"))
	}
}

func (cfg *apiConfig) handleValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type respVals struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error decoding parameters")
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "couldn't marshal JSON")
		return
	}

	cleanedBody := cleanedBody(params.Body)

	rv := respVals{
		CleanedBody: cleanedBody,
	}

	respondWithJSON(w, http.StatusOK, rv)
}

func cleanedBody(body string) string {
	splitBody := strings.Split(body, " ")
	cleanedBody := []string{}
	for _, sb := range splitBody {
		cleanedBody = append(cleanedBody, replaceBadWords(sb))
	}
	return strings.Join(cleanedBody, " ")
}

func replaceBadWords(word string) string {
	badWords := []string{"kerfuffle", "sharbert", "fornax"}
	for _, bw := range badWords {
		if strings.Contains(strings.ToLower(word), bw) {
			word = "****"
		}
	}
	return word
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}

	type respVals struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error decoding parameters")
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't create user")
	}

	rv := respVals{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}

	respondWithJSON(w, http.StatusCreated, rv)
}
