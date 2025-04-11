package main

import (
	"encoding/json"
	"net/http"
)

type AppErrorResp struct {
	Error string `json:"error"`
}

func RespondWithError(w http.ResponseWriter, code int, msg string) {
	errorResp := AppErrorResp{
		Error: msg,
	}
	errorRespBytes, err := json.Marshal(errorResp)
	if err != nil {
		w.WriteHeader(code)
	}
	w.WriteHeader(code)
	w.Write(errorRespBytes)
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.WriteHeader(code)
	bytes, err := json.Marshal(payload)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Error marshalling data")
		return
	}
	w.Write(bytes)
}
