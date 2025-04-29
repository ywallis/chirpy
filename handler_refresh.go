package main

import (
	"context"
	"net/http"
	"time"

	"github.com/ywallis/chirpy/internal/auth"
)

func (a *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(
			w,
			http.StatusBadRequest,
			"No token found in request",
			err,
		)
		return
	}

	dbToken, err := a.db.FetchRefreshToken(context.Background(), token)
	if err != nil {
		respondWithError(
			w,
			http.StatusUnauthorized,
			"Cannot find refresh token",
			err,
		)
		return
	}
	if dbToken.ExpiresAt.Before(time.Now()) {

		respondWithError(
			w,
			http.StatusUnauthorized,
			"Token has expired",
			err,
		)
		return
	}
	if dbToken.RevokedAt.Valid {

		respondWithError(
			w,
			http.StatusUnauthorized,
			"Token has been revoked",
			err,
		)
		return
	}
	newAccessToken, err := auth.MakeJWT(
		dbToken.UserID,
		a.JWTSecret,
		time.Duration(time.Minute * 60),
	)
	if err != nil {

		respondWithError(
			w,
			http.StatusInternalServerError,
			"Cannot generate new token",
			err,
		)
		return
	}
	response := AccessToken{
		Token: newAccessToken,
	}
	respondWithJSON(w, http.StatusOK, response)
}
