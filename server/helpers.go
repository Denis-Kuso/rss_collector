package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(payload)
	if err != nil {
		slog.Warn("JSON marshalling failed", "error", err, "payload", payload)
		// fallback
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(statusCode)
	w.Write(data)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type errorResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errorResponse{
		Error: msg,
	})
}

func (a *app) logError(r *http.Request, err error) {
	slog.Error("handlerError", "error", err, "req method", r.Method, "URL", r.URL)
}

func (app *app) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)

	msg := "OOPS! Ran into a problem and could not process your request"
	respondWithError(w, http.StatusInternalServerError, msg)
}
