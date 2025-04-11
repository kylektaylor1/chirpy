package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/kylektaylor1/chirpy/internal/auth"
	"github.com/kylektaylor1/chirpy/internal/database"
)

type params struct {
	Body   string    `json:"body"`
	UserId uuid.UUID `json:"user_id"`
}

type resp struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserId    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) HandlerCreateChirp(w http.ResponseWriter, r *http.Request) {
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
	dec := json.NewDecoder(r.Body)
	var reqBody params
	err := dec.Decode(&reqBody)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "error decoding req body")
		return
	}

	strlen := len(reqBody.Body)
	if strlen <= 140 {
		cleaned := cleanBody(reqBody.Body)

		chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
			Body:   cleaned,
			UserID: userId,
		})
		if err != nil {
			fmt.Printf("Error creating chirp: %v", err)
			RespondWithError(w, http.StatusInternalServerError, "error creating chirp")
			return
		}

		resp := resp{
			Id:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserId:    userId,
		}

		RespondWithJSON(w, http.StatusCreated, resp)
		return
	} else {

		RespondWithJSON(w, http.StatusBadRequest, AppErrorResp{
			Error: "Chirp is too long",
		})

		return
	}
}

func cleanBody(s string) string {
	badWords := []string{"kerfuffle", "sharbert", "fornax"}
	splitWords := strings.Split(s, " ")

	var cleanedWords []string
	for _, word := range splitWords {
		loweredWord := strings.ToLower(word)
		ok := slices.Contains(badWords, loweredWord)
		if ok {
			cleanedWords = append(cleanedWords, "****")
		} else {
			cleanedWords = append(cleanedWords, word)
		}
	}

	retStr := strings.Join(cleanedWords, " ")
	return retStr
}
