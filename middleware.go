package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/sam-maton/chirpy/internal/auth"
)

func (cfg *apiConfig) middlewareAuth(handler func(http.ResponseWriter, *http.Request, uuid.UUID)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
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

		handler(w, r, userId)
	}
}
