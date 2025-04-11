package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/kylektaylor1/chirpy/internal/auth"
)

type handlerUpgradeChirpyRedParams struct {
	Event string `json:"event"`
	Data  struct {
		UserId string `json:"user_id"`
	} `json:"data"`
}

func (cfg *apiConfig) HandlerUpgradeChirpyRed(w http.ResponseWriter, r *http.Request) {
	apiKey, apiErr := auth.GetAPIKey(r.Header)
	if apiErr != nil {
		RespondWithError(w, http.StatusUnauthorized, "no api key")
		return
	}

	if apiKey != cfg.polkaKey {
		RespondWithError(w, http.StatusUnauthorized, "invalid api key")
		return
	}

	dec := json.NewDecoder(r.Body)
	var reqBody handlerUpgradeChirpyRedParams
	err := dec.Decode(&reqBody)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "cannot parse params")
		return
	}

	if reqBody.Event != "user.upgraded" {
		RespondWithJSON(w, http.StatusNoContent, nil)
		return
	}

	upgradeErr := cfg.db.UpgradeUserById(r.Context(), uuid.MustParse(reqBody.Data.UserId))
	if upgradeErr != nil {
		RespondWithError(w, http.StatusNotFound, "user not upgraded")
		return
	}

	RespondWithJSON(w, http.StatusNoContent, nil)
}
