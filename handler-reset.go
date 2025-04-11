package main

import (
	"fmt"
	"net/http"
	"os"
)

func (cfg *apiConfig) HandlerResetUsers(w http.ResponseWriter, r *http.Request) {
	platform := os.Getenv("PLATFORM")
	if platform != "dev" {
		RespondWithError(w, http.StatusForbidden, "not platform = dev")
		return
	}
	err := cfg.db.DeleteUsers(r.Context())
	if err != nil {
		fmt.Printf("Error creating user: %v", err)
		RespondWithError(w, http.StatusInternalServerError, "error creating user")
		return
	}

	RespondWithJSON(w, http.StatusOK, nil)
}
