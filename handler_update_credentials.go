package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/ywallis/chirpy/internal/auth"
	"github.com/ywallis/chirpy/internal/database"
)

func (a *apiConfig) handlerUpdateCredentials(w http.ResponseWriter, r *http.Request) {

	type requestParams struct{
		Password string `json:"password"`
		Email string `json:"email"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(
			w,
			http.StatusUnauthorized,
			"No bearer token found",
			err,
		)
		return
	}
	requestItem := requestParams{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requestItem); err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"Couldn't decode parameters",
			err,
		)
		return
	}

	userID, err := auth.ValidateJWT(token, a.JWTSecret)
	if err != nil {
		respondWithError(
			w,
			http.StatusUnauthorized,
			"Token cannot be validated",
			err,
		)
		return
	}
	hashedPw, err := auth.HashPassword(requestItem.Password)
	updateParams := database.UpdateCredentialsParams{
		Email:    requestItem.Email,
		Password: hashedPw, 
		UpdatedAt: time.Now().UTC(),
		ID:  userID,
	}
	user, err := a.db.UpdateCredentials(context.Background(), updateParams)
	responseParams := User{
		Id:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
	}
	respondWithJSON(w, http.StatusOK, responseParams)


}
