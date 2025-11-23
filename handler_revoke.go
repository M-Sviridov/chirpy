package main

import (
	"net/http"

	"github.com/M-Sviridov/chirpy/internal/auth"
)

func (cfg *apiConfig) handleRevoke(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't get bearer token")
		return
	}

	user, err := cfg.db.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "couldn't get user from refresh token")
		return
	}

	if user.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "refresh token has already expired")
		return
	}

	if err := cfg.db.RevokeRefreshToken(r.Context(), refreshToken); err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't revoke frefreshtoken")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
