package main

import (
	"net/http"

	"github.com/kylektaylor1/chirpy/internal/auth"
)

func (cfg *apiConfig) HandlerRevoke(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "error getting bearer token")
		return
	}

	updateErr := cfg.db.SetRefreshTokenRevoked(r.Context(), token)
	if updateErr != nil {
		RespondWithError(w, http.StatusInternalServerError, "refresh token not updated")
		return
	}

	RespondWithJSON(w, http.StatusNoContent, struct{}{})
}
