package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRoute(t *testing.T) {
	a := app{}
	t.Run("hit non-existing endpoint", func(t *testing.T) {
		endpoint := "/gilmore"
		req, _ := http.NewRequest("GET", endpoint, nil)
		rr := httptest.NewRecorder()
		h := a.setupRoutes() // TODO: is this best way to do it
		h.ServeHTTP(rr, req)
		assertStatus(t, rr.Code, http.StatusNotFound)
	})

	t.Run("wrong method", func(t *testing.T) {
		endpoint := "/v1/users"
		req, _ := http.NewRequest("PATCH", endpoint, nil)
		rr := httptest.NewRecorder()
		h := a.setupRoutes() // TODO: is this best way to do it
		h.ServeHTTP(rr, req)
		// TODO body could be better
		assertStatus(t, rr.Code, http.StatusMethodNotAllowed)

	})
}
