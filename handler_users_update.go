package main

import (
	"encoding/json"
	"net/http"

	"github.com/M-Sviridov/chirpy/internal/auth"
	"github.com/M-Sviridov/chirpy/internal/database"
)

func (cfg *apiConfig) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type respVals struct {
		Email string `json:"email"`
	}

	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "couldn't get bearer token")
		return
	}

	userID, err := auth.ValidateJWT(accessToken, cfg.tokenSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "JWT token is invalid")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error decoding parameters")
		return
	}

	hash, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't hash password")
		return
	}

	userParams := database.UpdateUserParams{
		Email:          params.Email,
		HashedPassword: hash,
		ID:             userID,
	}

	rv := respVals{
		Email: userParams.Email,
	}

	if err := cfg.db.UpdateUser(r.Context(), userParams); err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't update user details in DB")
	}

	respondWithJSON(w, http.StatusOK, rv)
}
