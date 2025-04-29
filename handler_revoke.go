package main

import (
	"context"
	"net/http"
	"time"

	"github.com/ywallis/chirpy/internal/auth"
	"github.com/ywallis/chirpy/internal/database"
)

func (a *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {

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
	revokeParams := database.RevokeRefreshTokenParams{
		UpdatedAt: time.Now(),
		Token:     token,
	}
	if err := a.db.RevokeRefreshToken(context.Background(), revokeParams); err != nil {

		respondWithError(
			w,
			http.StatusInternalServerError,
			"Error revoking token",
			err,
		)
		return
	}
	respondWithJSON(w, http.StatusNoContent, nil)
}
