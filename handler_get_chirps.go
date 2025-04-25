package main

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

func (a *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {

	chirps, err := a.db.GetAllChirps(context.Background())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not fetch chirps", err)
	}

	output := []Chirp{}
	for _, chirp := range chirps {
		item := Chirp{
			Id:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		}
		output = append(output, item)
	}
	respondWithJSON(w, http.StatusOK, output)
}

func (a *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")

	uuid, err := uuid.Parse(id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "ID does not follow UUID standards", err)
	}

	chirp, err := a.db.GetChirp(context.Background(), uuid)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Could not find chirp", err)
	}
	output := Chirp{
		Id:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}
	respondWithJSON(w, http.StatusOK, output)
}
