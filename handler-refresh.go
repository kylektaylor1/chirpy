package main

import (
	"net/http"
	"time"

	"github.com/kylektaylor1/chirpy/internal/auth"
)

type handlerRefreshResp struct {
	Token string `json:"token"`
}

func (cfg *apiConfig) HandlerRefresh(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "error getting bearer token")
		return
	}

	rToken, err := cfg.db.GetRefreshTokenByToken(r.Context(), token)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "refresh token not in db")
		return
	}

	if time.Now().Compare(rToken.ExpiresAt) == 1 {
		RespondWithError(w, http.StatusUnauthorized, "expired refresh token")
		return
	}

	if rToken.RevokedAt.Valid {
		RespondWithError(w, http.StatusUnauthorized, "token revoked")
		return
	}

	jwt, err := auth.MakeJWT(rToken.UserID, cfg.jwtSecret, time.Hour)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "cannot gen jwt")
		return
	}

	RespondWithJSON(w, http.StatusOK, handlerRefreshResp{
		Token: jwt,
	})
}
