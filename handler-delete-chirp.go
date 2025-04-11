package main

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/kylektaylor1/chirpy/internal/auth"
)

func (cfg *apiConfig) HandlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	token, getTokenErr := auth.GetBearerToken(r.Header)
	if getTokenErr != nil {
		RespondWithError(w, http.StatusUnauthorized, "no auth token")
		return
	}
	userId, validErr := auth.ValidateJWT(token, cfg.jwtSecret)
	if validErr != nil {
		fmt.Printf("validErr: %v\n", validErr)
		RespondWithError(w, http.StatusUnauthorized, "invalid token")
		return
	}

	path := r.PathValue("chirpID")
	chirpId, err := uuid.Parse(path)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "invalid uuid")
		return
	}

	chirp, err := cfg.db.GetChirpById(r.Context(), chirpId)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, "cannot query chirp")
		return
	}

	if chirp.UserID != userId {
		RespondWithError(w, http.StatusForbidden, "forbidden")
		return
	}

	deleteErr := cfg.db.DeleteChirpById(r.Context(), chirpId)
	if deleteErr != nil {
		RespondWithError(w, http.StatusInternalServerError, "cannot delete chirp")
		return
	}

	RespondWithJSON(w, http.StatusNoContent, nil)
}
