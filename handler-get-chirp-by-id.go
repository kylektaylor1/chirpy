package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type getOneChirpResp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) HandlerGetChirpById(w http.ResponseWriter, r *http.Request) {
	path := r.PathValue("chirpID")
	chirpId, err := uuid.Parse(path)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "invalid uuid")
	}
	chirp, err := cfg.db.GetChirpById(r.Context(), chirpId)
	if err != nil {
		fmt.Printf("Error getting chirp: %v", err)
		RespondWithError(w, http.StatusNotFound, "chirp not found")
		return
	}

	resp := getOneChirpResp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}

	RespondWithJSON(w, http.StatusOK, resp)
}
