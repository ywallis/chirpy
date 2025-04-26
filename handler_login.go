package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/ywallis/chirpy/internal/auth"
)

func (a *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {

	type requestParams struct {
		Email         string `json:"email"`
		Password      string `json:"password"`
		ExpirationSec int    `json:"expires_in_seconds,omitempty"`
	}

	decoder := json.NewDecoder(r.Body)
	params := requestParams{}
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

	if params.ExpirationSec < 1 || params.ExpirationSec > 3600 {
		params.ExpirationSec = 3600
	}
	user, err := a.db.GetUserByEmail(context.Background(), params.Email)
	if err != nil {
		respondWithError(
			w,
			http.StatusUnauthorized,
			"Incorrect email or password",
			err,
		)
		return
	}
	if err := auth.CheckPasswordHash(params.Password, user.Password); err != nil {
		respondWithError(
			w,
			http.StatusUnauthorized,
			"Incorrect email or password",
			err,
		)
		return
	}
	token, err := auth.MakeJWT(user.ID, a.JWTSecret, time.Duration(params.ExpirationSec)*time.Second)
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"Could not generate token",
			err,
		)
		return
	}
	output := User{
		Id:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Token:     token,
	}
	respondWithJSON(w, http.StatusOK, output)
}
