package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ywallis/chirpy/internal/auth"
	"github.com/ywallis/chirpy/internal/database"
)

func (a *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {

	type inputParameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := inputParameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if params.Password == "" {
		respondWithError(w, http.StatusBadRequest, "Password cannot be empty", err)
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error hashing password", err)
	}
	queryParams := database.CreateUserParams{
		Email:    params.Email,
		Password: hashedPassword,
	}
	user, err := a.db.CreateUser(context.Background(), queryParams)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Could not create user", err)
		return
	}
	responseValues := User{
		Id:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	}
	respondWithJSON(w, http.StatusCreated, responseValues)
}
