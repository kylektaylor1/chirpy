package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/kylektaylor1/chirpy/internal/auth"
	"github.com/kylektaylor1/chirpy/internal/database"
)

type userLoginParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type userLoginResp struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
}

func (cfg *apiConfig) HandlerUserLogin(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	var params userLoginParams
	err := dec.Decode(&params)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "error decoding params")
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	authErr := auth.CheckPasswordHash(user.HashedPassword, params.Password)
	if authErr != nil {
		RespondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	jwt, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "cannot create jwt")
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "cannot create refresh token")
		return
	}

	refreshRecord, err := cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 60),
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "cannot write refresh token to db")
		return
	}

	RespondWithJSON(w, http.StatusOK, userLoginResp{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        jwt,
		RefreshToken: refreshRecord.Token,
		IsChirpyRed:  user.IsChirpyRed,
	})
}
