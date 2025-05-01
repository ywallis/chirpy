package main

import (
	"context"
	"net/http"
	"sort"

	"github.com/google/uuid"
	"github.com/ywallis/chirpy/internal/auth"
	"github.com/ywallis/chirpy/internal/database"
)

func (a *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {

	var chirps []database.Chirp
	var err error
	var sortOrder string
	sortOrderParam := r.URL.Query().Get("sort")
	switch sortOrderParam {
	case "asc": 
		sortOrder = ""
	case "":
	sortOrder = ""
	case "desc":
	sortOrder = "desc"
	default:

			respondWithError(
				w,
				http.StatusBadRequest,
				"Invalid sort parameter",
				err,
			)
			return
}
	author := r.URL.Query().Get("author_id")
	// Get chirps from specific author if requested, all otherwise
	if author != "" {
		userId, parseErr := uuid.Parse(author)
		if parseErr != nil {

			respondWithError(
				w,
				http.StatusBadRequest,
				"Invalid author ID",
				err,
			)
			return
		}
		chirps, err = a.db.GetChirpsFromUser(context.Background(), userId)
	} else {
		chirps, err = a.db.GetAllChirps(context.Background())
	}
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"Could not fetch chirps",
			err,
		)
		return
	}

	output := []Chirp{}
	for _, chirp := range chirps {
		item := Chirp{
			Id:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		}
		output = append(output, item)
	}
	if sortOrder == "desc" {
		sort.Slice(output, func(i, j int) bool {
			return output[i].CreatedAt.After(output[j].CreatedAt)
		})
	}
	respondWithJSON(w, http.StatusOK, output)
}

func (a *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")

	uuid, err := uuid.Parse(id)
	if err != nil {
		respondWithError(
			w,
			http.StatusBadRequest,
			"ID does not follow UUID standards",
			err,
		)
		return
	}

	chirp, err := a.db.GetChirp(context.Background(), uuid)
	if err != nil {
		respondWithError(
			w,
			http.StatusNotFound,
			"Could not find chirp",
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
	respondWithJSON(w, http.StatusOK, output)
}

func (a *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(
			w,
			http.StatusUnauthorized,
			"No token found in header",
			err,
		)
		return
	}
	userId, err := auth.ValidateJWT(token, a.JWTSecret)
	if err != nil {
		respondWithError(
			w,
			http.StatusUnauthorized,
			"Token cannot be validated",
			err,
		)
		return
	}
	chirpId := r.PathValue("id")
	uuid, err := uuid.Parse(chirpId)
	if err != nil {
		respondWithError(
			w,
			http.StatusBadRequest,
			"ID does not follow UUID standards",
			err,
		)
		return
	}

	chirp, err := a.db.GetChirp(context.Background(), uuid)
	if err != nil {
		respondWithError(
			w,
			http.StatusNotFound,
			"Could not find chirp",
			err,
		)
		return
	}
	if chirp.UserID != userId {
		respondWithError(
			w,
			http.StatusForbidden,
			"User not authorized",
			err,
		)
		return
	}
	if err := a.db.DeleteChirp(context.Background(), chirp.ID); err != nil {

		respondWithError(
			w,
			http.StatusInternalServerError,
			"Chirp cannot be deleted",
			err,
		)
		return
	}
	respondWithJSON(w, http.StatusNoContent, nil)

}
