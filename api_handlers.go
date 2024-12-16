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
		Email    string `json:"email"`
		Password string `json:"password"`
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

	accessToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, "Couldn't create JWT", http.StatusInternalServerError, err)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, "Couldn't create refresh token", http.StatusInternalServerError, err)
		return
	}

	refreshTokenParams := database.CreateRefreshTokenParams{
		Token:     refreshToken,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		ExpiresAt: (time.Now().Add(time.Hour * 1440)),
		UserID:    user.ID,
	}
	cfg.db.CreateRefreshToken(r.Context(), refreshTokenParams)

	respondWithJson(w, 200, response{
		User: database.User{
			ID:    user.ID,
			Email: user.Email,
		},
		Token:        accessToken,
		RefreshToken: refreshToken,
	},
	)

}

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	headerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, "There was no refresh token in the header", http.StatusUnauthorized, err)
		return
	}

	refreshToken, err := cfg.db.GetRefreshTokenByToken(r.Context(), headerToken)
	if err != nil {
		respondWithError(w, "There was no refresh token", http.StatusUnauthorized, err)
	}

	if time.Now().After(refreshToken.ExpiresAt) {
		respondWithError(w, "The refresh token has expired", http.StatusUnauthorized, nil)
	}

	user, err := cfg.db.GetUserByRefreshToken(r.Context(), refreshToken.Token)
	if err != nil {
		respondWithError(w, "There was an error getting the user by refresh token", http.StatusInternalServerError, err)
	}

	newToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, "There was an error creating the new JWT", http.StatusInternalServerError, err)
	}

	// Return refresh token
	respondWithJson(w, 200, struct{ token string }{
		token: newToken,
	})
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, "There was no auth token in the header", http.StatusUnauthorized, err)
		return
	}

	userId, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, "Invalid auth token", http.StatusUnauthorized, err)
		return
	}

	type params struct {
		Body string `json:"body"`
	}
	p := params{}
	decoder := json.NewDecoder(r.Body)

	err = decoder.Decode(&p)
	if err != nil {
		respondWithError(w, paramsDecodeError, 400, err)
		return
	}

	createParams := database.CreateChirpParams{
		Body:   p.Body,
		UserID: userId,
	}

	chirp, err := cfg.db.CreateChirp(r.Context(), createParams)
	if err != nil {
		respondWithError(w, "There was an error creating the chirp", 400, err)
		return
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
