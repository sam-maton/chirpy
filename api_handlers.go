package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/sam-maton/chirpy/internal/database"
)

const paramsDecodeError = "There was an error decoding the params"

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
		respondWithError(w, paramsDecodeError, http.StatusInternalServerError, err)
		return
	}

	if len(p.Body) > 140 {
		respondWithError(w, "Chirp is too long", 400, nil)
		return
	}

	clean := cleanChirp(p.Body)

	respondWithJson(w, 200, validParams{CleanedBody: clean})
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type params struct {
		Email string `json:"email"`
	}
	decoder := json.NewDecoder(r.Body)
	p := params{}
	err := decoder.Decode(&p)

	if err != nil {
		respondWithError(w, paramsDecodeError, http.StatusInternalServerError, err)
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), p.Email)

	if err != nil {
		respondWithError(w, "There was an error creating the user", http.StatusInternalServerError, err)
		return
	}

	respondWithJson(w, 201, user)
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	type params struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}
	p := params{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&p)

	if err != nil {
		respondWithError(w, paramsDecodeError, 400, err)
	}

	createParams := database.CreateChirpParams{
		Body:   p.Body,
		UserID: p.UserID,
	}

	chirp, err := cfg.db.CreateChirp(r.Context(), createParams)

	if err != nil {
		respondWithError(w, "There was an error creating the chirp", 400, err)
	}

	respondWithJson(w, 201, chirp)
}

func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {

	chirps, err := cfg.db.GetChirps(r.Context())

	if err != nil {
		respondWithError(w, "There was an error getting all the chirps", 400, err)
	}

	respondWithJson(w, 200, chirps)
}
