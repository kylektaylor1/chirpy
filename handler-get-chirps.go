package main

import (
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/kylektaylor1/chirpy/internal/database"
)

type getChirpsStruct struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) HandlerGetChirps(w http.ResponseWriter, r *http.Request) {
	s := r.URL.Query().Get("author_id")
	sortParam := r.URL.Query().Get("sort")
	var chirps []database.Chirp
	if s == "" {
		allChirps, err := cfg.db.GetAllChirps(r.Context())
		if err != nil {
			fmt.Printf("Error getting chirps: %v", err)
			RespondWithError(w, http.StatusInternalServerError, "error getting chirps")
			return
		}
		chirps = append(chirps, allChirps...)
	} else {
		allChirps, err := cfg.db.GetAllChirpsByUserId(r.Context(), uuid.MustParse(s))
		if err != nil {
			fmt.Printf("Error getting chirps: %v", err)
			RespondWithError(w, http.StatusInternalServerError, "error getting chirps")
			return
		}
		chirps = append(chirps, allChirps...)
	}

	var resp []getChirpsStruct
	for _, chirp := range chirps {
		resp = append(resp, getChirpsStruct{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}
	var sortValue int
	if sortParam == "asc" {
		sortValue = -1
	} else {
		sortValue = 1
	}
	sort.Slice(resp, func(i, j int) bool { return resp[i].CreatedAt.Compare(resp[j].CreatedAt) == sortValue })

	RespondWithJSON(w, http.StatusOK, resp)
}
