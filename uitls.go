package main

import (
	"sort"

	"github.com/sam-maton/chirpy/internal/database"
)

func sortChirpsDesc(ascChirps []database.Chirp) []database.Chirp {

	sort.Slice(ascChirps, func(i, j int) bool {
		return ascChirps[i].CreatedAt.After(ascChirps[j].CreatedAt)
	})

	return ascChirps
}
