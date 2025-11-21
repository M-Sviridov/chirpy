package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/M-Sviridov/chirpy/internal/auth"
	"github.com/M-Sviridov/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
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
	}

	hash, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't hash password")
	}

	userParams := database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hash,
	}

	user, err := cfg.db.CreateUser(r.Context(), userParams)
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
