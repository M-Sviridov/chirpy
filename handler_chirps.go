package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/M-Sviridov/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handleCreateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	type respVals struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error decoding parameters")
		return
	}

	cleanedBody, err := cfg.validateChirp(params.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't validate chirp")
	}

	chirpParams := database.CreateChirpParams{
		Body:   cleanedBody,
		UserID: params.UserID,
	}

	chirp, err := cfg.db.CreateChirp(r.Context(), chirpParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't create chirp")
		log.Fatalf("couldn't create chirp: %s", err)
	}

	rv := respVals{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    params.UserID,
	}

	respondWithJSON(w, http.StatusCreated, rv)
}

func (cfg *apiConfig) validateChirp(body string) (string, error) {
	const maxChirpLength = 140
	if len(body) > maxChirpLength {
		return "", errors.New("chirp length cannot exceed 140 characters")
	}

	cleanedBody := cleanedBody(body)

	return cleanedBody, nil
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
