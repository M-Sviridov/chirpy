package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

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
