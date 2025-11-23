package main

import (
	"net/http"
	"time"

	"github.com/M-Sviridov/chirpy/internal/auth"
)

func (cfg *apiConfig) handleRefresh(w http.ResponseWriter, r *http.Request) {
	type respVals struct {
		Token string `json:"token"`
	}

	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "couldn't get bearer token")
		return
	}

	user, err := cfg.db.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "couldn't get user from refresh token")
		return
	}

	if user.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "refresh token has expired")
		return
	}

	accessToken, err := auth.MakeJWT(user.UserID, cfg.tokenSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't create access token")
		return

	}

	rv := respVals{
		Token: accessToken,
	}

	respondWithJSON(w, http.StatusOK, rv)
}
