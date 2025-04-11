package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/kylektaylor1/chirpy/internal/auth"
	"github.com/kylektaylor1/chirpy/internal/database"
)

func (cfg *apiConfig) HandlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type params struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type resp struct {
		Id          uuid.UUID `json:"id"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
		Email       string    `json:"email"`
		IsChirpyRed bool      `json:"is_chirpy_red"`
	}
	dec := json.NewDecoder(r.Body)
	var reqBody params
	err := dec.Decode(&reqBody)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "error decoding req body")
		return
	}

	hpw, err := auth.HashPassword(reqBody.Password)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "error hashing pw")
	}

	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:          reqBody.Email,
		HashedPassword: hpw,
	})
	if err != nil {
		fmt.Printf("Error creating user: %v", err)
		RespondWithError(w, http.StatusInternalServerError, "error creating user")
		return
	}

	RespondWithJSON(w, http.StatusCreated, resp{
		Id:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	})
}
