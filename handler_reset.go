package main

import "net/http"

func (cfg *apiConfig) handleReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, http.StatusForbidden, "cannot run reset on non dev platform")
		return
	}

	cfg.fileserverHits.Store(0)
	err := cfg.db.DeleteUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't delete all users")
		return
	}

	respondWithJSON(w, http.StatusOK, "reset hits to 0 and delete all users in database")
}
