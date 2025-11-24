package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handleWebhooks(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	// type respVals struct {
	// 	ID        uuid.UUID `json:"id"`
	// 	CreatedAt time.Time `json:"created_at"`
	// 	UpdatedAt time.Time `json:"updated_at"`
	// 	Email     string    `json:"email"`
	// }

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error decoding parameters")
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	userID, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't parse user ID")
		return
	}

	if err := cfg.db.UpgradeUserToChirpyRed(r.Context(), userID); err != nil {
		respondWithError(w, http.StatusNotFound, "couldn't upgrade user to chirpy red")
		return
	}

	w.WriteHeader(http.StatusNoContent)

	// rv := respVals{
	// 	ID:        user.ID,
	// 	CreatedAt: user.CreatedAt,
	// 	UpdatedAt: user.UpdatedAt,
	// 	Email:     user.Email,
	// }

	// respondWithJSON(w, http.StatusCreated, rv)
}
