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

type handlerPutUsersParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type handlerPutUsersResp struct {
	Id          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Email       string    `json:"email"`
	IsChirpyRed bool      `json:"is_chirpy_red"`
}

func (cfg *apiConfig) HandlerPutUsers(w http.ResponseWriter, r *http.Request) {
	jwt, tokenErr := auth.GetBearerToken(r.Header)
	if tokenErr != nil {
		RespondWithError(w, http.StatusUnauthorized, "error getting bearer token")
		return
	}

	userId, err := auth.ValidateJWT(jwt, cfg.jwtSecret)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "invalid bearer token")
		return
	}

	dec := json.NewDecoder(r.Body)
	var reqBody handlerPutUsersParams
	decErr := dec.Decode(&reqBody)
	if decErr != nil {
		RespondWithError(w, http.StatusInternalServerError, "error decoding req body")
		return
	}

	fmt.Printf("useremail: %v\n", reqBody.Email)
	user, err := cfg.db.GetUserById(r.Context(), userId)
	if err != nil {
		fmt.Printf("Error getting user: %v", err)
		RespondWithError(w, http.StatusInternalServerError, "error getting user")
		return
	}

	if user.ID != userId {
		RespondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	newHashedPw, err := auth.HashPassword(reqBody.Password)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "error hashing pw")
		return
	}

	updateErr := cfg.db.UpdateUser(r.Context(), database.UpdateUserParams{
		Email:          reqBody.Email,
		HashedPassword: newHashedPw,
		ID:             userId,
	})
	if updateErr != nil {
		RespondWithError(w, http.StatusInternalServerError, "error updating user record")
		return
	}

	RespondWithJSON(w, http.StatusOK, handlerPutUsersResp{
		Id:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       reqBody.Email,
		IsChirpyRed: user.IsChirpyRed,
	})
}
