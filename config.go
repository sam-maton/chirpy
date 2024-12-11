package main

import (
	"database/sql"
	"log"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	"github.com/sam-maton/chirpy/internal/database"
)

func setupConfig() apiConfig {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbURL)

	if err != nil {
		log.Printf("There was an error connecting to the db: %s", err)
		os.Exit(1)
	}

	dbQueries := database.New(db)

	return apiConfig{
		fileServerHits: atomic.Int32{},
		db:             dbQueries,
	}
}

type apiConfig struct {
	fileServerHits atomic.Int32
	db             *database.Queries
}
