package main

import (
	"context"
	"encoding/json"
	"net/http"
)

func (a *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {

	type inputParameters struct {
		Email string `json:"email"`
	}


	decoder := json.NewDecoder(r.Body)
	params := inputParameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	user, err := a.db.CreateUser(context.Background(), params.Email)
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
	respondWithJSON(w, http.StatusCreated, responseValues)
}
