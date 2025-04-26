package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ywallis/chirpy/internal/auth"
	"github.com/ywallis/chirpy/internal/database"
)

func (a *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"Couldn't decode parameters",
			err,
		)
		return
	}
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(
			w,
			http.StatusUnauthorized,
			"Cannot find token",
			err,
		)
		return
	}

	userUuid, err := auth.ValidateJWT(token, a.JWTSecret)
	if err != nil {
		respondWithError(
			w,
			http.StatusUnauthorized,
			"Token not valid",
			err,
		)
		return
	}
	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(
			w,
			http.StatusBadRequest,
			"Chirp is too long",
			nil,
		)
		return
	}

	createChirpParams := database.CreateChirpParams{
		Body:   params.Body,
		UserID: userUuid,
	}
	chirp, err := a.db.CreateChirp(context.Background(), createChirpParams)
	if err != nil {
		respondWithError(
			w,
			http.StatusBadRequest,
			"Could not create Chirp",
			err,
		)
		return
	}

	output := Chirp{
		Id:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}
	respondWithJSON(w, http.StatusCreated, output)
}
