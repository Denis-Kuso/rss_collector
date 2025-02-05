package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/google/uuid"
)

func GetUserIDFromContext(r *http.Request) (uuid.UUID, bool) {
	uIDTemp := r.Context().Value(userIDctx)
	if uIDTemp == nil || (uIDTemp == uuid.UUID{}) {
		return uuid.UUID{}, false
	}

	userID, ok := uIDTemp.(uuid.UUID)
	if !ok {
		return uuid.UUID{}, false
	}

	return userID, true
}

func readJSON(r *http.Request, target interface{}) error {
	// io.EOF is not returned on empty body
	if r.Body == nil {
		return errors.New("body must not be empty")
	}

	dc := json.NewDecoder(r.Body)
	dc.DisallowUnknownFields()
	err := dc.Decode(target)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		switch {
		// see https://github.com/golang/go/issues/25956
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("poorly-formed JSON")
		// helpful message for well-intentioned people
		case errors.As(err, &syntaxError):
			return fmt.Errorf("poorly-formed JSON (at position %d)", syntaxError.Offset)
		// on wrong target type
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("wrong JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("wrong JSON type (at position %d)", unmarshalTypeError.Offset)
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")
		//still an open issue #29035 not making it into a type
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body has unsupported key %s", fieldName)
		// this error should never happen
		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		default:
			return err
		}
	}
	return nil
}

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

func (a *app) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	err = fmt.Errorf("%w: %s", err, string(debug.Stack()))
	a.logError(r, err)

	msg := "OOPS! Ran into a problem and could not process your request"
	respondWithError(w, http.StatusInternalServerError, msg)
}
