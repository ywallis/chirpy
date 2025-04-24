package main

import (
	"context"
	"encoding/json"
	"net/http"
)

func (a *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {

	type inputParameters struct {
		Body string `json:"body"`
	}


	decoder := json.NewDecoder(r.Body)
	params := inputParameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	user, err := a.db.CreateUser(context.Background(), params.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Could not create user", err)
		return
	}
	responseValues := User{
		Id: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
	}
	respondWithJSON(w, http.StatusOK, responseValues)
}
