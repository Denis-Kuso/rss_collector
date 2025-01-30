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

func readJSON(r *http.Request, target interface{}) error { // Decode the request body into the target destination.
	//TODO does the io.EOF solve this problem
	if r.Body == nil {
		return errors.New("body must not be empty")
	}
	//////////////////maxBytes := 1_048_576 // todo how big of a request
	//r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dc := json.NewDecoder(r.Body)
	dc.DisallowUnknownFields()
	err := dc.Decode(target)
	if err != nil {
		// If there is an error during decoding, start the triage...
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
		// Likewise, catch any *json.UnmarshalTypeError errors. These occur when the
		// JSON value is the wrong type for the target destination. If the error relates // to a specific field, then we include that in our error message to make it
		// easier for the client to debug.
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("wrong JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("wrong JSON type (at position %d)", unmarshalTypeError.Offset)
		// An io.EOF error will be returned by Decode() if the request body is empty. We // check for this with errors.Is() and return a plain-english error message
		// instead.
		// this could panic the decoder
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")
		//still an open issue #29035 not making it into a type
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body has unsupported key %s", fieldName)
		// A json.InvalidUnmarshalError error will be returned if we pass a non-nil
		// pointer to Decode(). We catch this and panic, rather than returning an error // to our handler. At the end of this chapter we'll talk about panicking
		// versus returning errors, and discuss why it's an appropriate thing to do in // this specific situation.
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
