package main

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"
)

func respondWithError(w http.ResponseWriter, message string, code int, err error) {
	if err != nil {
		// log.Println("Print error message")
		log.Println(err)
	}

	type errorParams struct {
		Error string `json:"error"`
	}

	respondWithJson(w, code, errorParams{Error: message})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	dat, err := json.Marshal(payload)

	if err != nil {
		log.Printf("There was an error marshaling the json: %s", err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	w.Write(dat)
}

func cleanChirp(message string) string {
	badWords := []string{"kerfuffle", "sharbert", "fornax"}

	splitMessage := strings.Split(message, " ")

	for i, v := range splitMessage {
		if slices.Contains(badWords, strings.ToLower(v)) {
			splitMessage[i] = "****"
		}
	}
	return strings.Join(splitMessage, " ")
}
