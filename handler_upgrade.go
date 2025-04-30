package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/ywallis/chirpy/internal/auth"
)

func (a *apiConfig) handlerUpgrade(w http.ResponseWriter, r *http.Request) {

	type data struct {
		UserID uuid.UUID `json:"user_id"`
	}
	type request struct {
		Event string `json:"event"`
		Data  data   `json:"data"`
	}
	apiKey, err := auth.GetApiKey(r.Header)
	if err != nil {
		respondWithError(
			w,
			http.StatusUnauthorized,
			"Couldn't find any api key",
			err,
		)
		return
	}
	if apiKey != a.polkaKey {

		respondWithError(
			w,
			http.StatusUnauthorized,
			"Invalid key",
			err,
		)
		return
	}
	requestBody := request{}
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err := decoder.Decode(&requestBody); err != nil {

		respondWithError(
			w,
			http.StatusInternalServerError,
			"Couldn't decode parameters",
			err,
		)
		return
	}

	if requestBody.Event != "user.upgraded" {

		respondWithJSON(w, http.StatusNoContent, nil)
		return
	}

	if err := a.db.MakePremium(
		context.Background(),
		requestBody.Data.UserID,
	); err != nil {

		respondWithError(
			w,
			http.StatusNotFound,
			"Couldn't find user",
			err,
		)
		return
	}
	respondWithJSON(w, http.StatusNoContent, nil)
}
