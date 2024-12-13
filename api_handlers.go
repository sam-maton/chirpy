package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sam-maton/chirpy/internal/auth"
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
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(r.Body)
	p := params{}
	err := decoder.Decode(&p)

	if err != nil {
		respondWithError(w, paramsDecodeError, http.StatusInternalServerError, err)
		return
	}

	hashedPW, err := auth.HashPassword(p.Password)

	if err != nil {
		respondWithError(w, "There was an error hashing the password", http.StatusInternalServerError, err)
		return
	}

	userParams := database.CreateUserParams{
		Email:          p.Email,
		HashedPassword: hashedPW,
	}

	user, err := cfg.db.CreateUser(r.Context(), userParams)

	if err != nil {
		respondWithError(w, "There was an error creating the user", http.StatusInternalServerError, err)
		return
	}

	respondWithJson(w, 201, user)
}

func (cfg *apiConfig) handlerLoginUser(w http.ResponseWriter, r *http.Request) {
	type params struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}

	type response struct {
		User         database.User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	p := params{}
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&p)
	if err != nil {
		respondWithError(w, paramsDecodeError, 400, err)
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), p.Email)
	if err != nil {
		respondWithError(w, "No user exists with that email address", 400, err)
	}

	err = auth.CheckPasswordHash(p.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, "Wrong password", 401, err)
		return
	}

	expireTime := time.Hour

	if p.ExpiresInSeconds > 0 && p.ExpiresInSeconds < 3600 {
		expireTime = time.Duration(p.ExpiresInSeconds) * time.Second
	}

	accessToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret, expireTime)
	if err != nil {
		respondWithError(w, "Couldn't create JWT", http.StatusInternalServerError, err)
		return
	}

	respondWithJson(w, 200, response{
		User: database.User{
			ID:    user.ID,
			Email: user.Email,
		},
		Token: accessToken,
	},
	)

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

func (cfg *apiConfig) handlerGetOneChirp(w http.ResponseWriter, r *http.Request) {

	pathID := r.PathValue("id")

	id, err := uuid.Parse(pathID)

	if err != nil {
		respondWithError(w, "Not a valid ID", 400, err)
		return
	}

	chirp, err := cfg.db.GetChirpByID(r.Context(), id)

	if err != nil {
		respondWithError(w, "There was an error getting the chirp by ID", 400, err)
		return
	}

	respondWithJson(w, 200, chirp)
}
