package main

import (
	"context"
	"fmt"
	"net/http"
)

func (a *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {

	if a.platform != "dev" {
		respondWithError(w, http.StatusForbidden, "Access forbidden", fmt.Errorf("Only allowed in development"))
	}

	a.db.DeleteAllUsers(context.Background())
	w.WriteHeader(200)
	w.Write([]byte("All gone"))
}
