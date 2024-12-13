package main

import (
	"database/sql"
	"log"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	"github.com/sam-maton/chirpy/internal/database"
)

type apiConfig struct {
	fileServerHits atomic.Int32
	db             *database.Queries
	jwtSecret      string
}

func setupConfig() apiConfig {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	secret := os.Getenv("SECRET")

	if dbURL == "" {
		log.Fatal("DB_URL environment variable must be set")
	}

	if secret == "" {
		log.Fatal("SECRET environment variable must be set")
	}

	db, err := sql.Open("postgres", dbURL)

	if err != nil {
		log.Printf("There was an error connecting to the db: %s", err)
		os.Exit(1)
	}

	dbQueries := database.New(db)

	return apiConfig{
		fileServerHits: atomic.Int32{},
		db:             dbQueries,
		jwtSecret:      secret,
	}
}
