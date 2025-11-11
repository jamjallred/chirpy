package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"slices"
	"strings"
)

func handlerReadiness(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func handlerValidateChirp(w http.ResponseWriter, req *http.Request) {

	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		log.Printf("Error decoding parameters: %s\n", err)
		respondWithError(w, 500, "Something went wrong")
		return
	}

	if len(params.Body) > 140 {
		log.Printf("Chirp is over 140 characters\n")
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	badWords := []string{"kerfuffle", "sharbert", "fornax"}

	chirpWords := strings.Split(params.Body, " ")
	for idx, word := range chirpWords {
		if slices.Contains(badWords, strings.ToLower(word)) {
			chirpWords[idx] = "****"
		}
	}

	type returnVal struct {
		Cleaned_body string `json:"cleaned_body"`
	}

	respBody := returnVal{
		Cleaned_body: strings.Join(chirpWords, " "),
	}

	respondWithJSON(w, http.StatusOK, respBody)

}

func respondWithError(w http.ResponseWriter, code int, msg string) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	type errVal struct {
		Error string `json:"error"`
	}

	returnVal := errVal{}

	dat, err := json.Marshal(returnVal)
	if err != nil {
		log.Printf("Error marshalling JSON: %s\n", err)
		return errors.New("error marshalling JSON")
	}

	w.Write(dat)

	return nil
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s\n", err)
		return errors.New("error marshalling JSON")
	}

	w.Write(dat)
	return nil
}
