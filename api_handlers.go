package main

import (
	"encoding/json"
	"net/http"
)

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type params struct {
		Body string `json:"body"`
	}

	type validParams struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	p := params{}

	err := decoder.Decode(&p)

	if err != nil {
		respondWithError(w, "There was an error decoding the params.", http.StatusInternalServerError, err)
		return
	}

	if len(p.Body) > 140 {
		respondWithError(w, "Chirp is too long", 400, nil)
		return
	}

	clean := cleanChirp(p.Body)

	respondWithJson(w, 200, validParams{CleanedBody: clean})
}

func handlerCreateUser(w http.ResponseWriter, r *http.Request) {

}
