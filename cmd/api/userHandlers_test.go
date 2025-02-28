package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Denis-Kuso/rss_collector/internal/storage"
	"github.com/google/uuid"
)

type MockUserStore struct {
	num   int // perhaps replace with a map
	users map[uuid.UUID]storage.User
}

var m MockUserStore

func TestGetUser(t *testing.T) {
	endpoint := "/v1/users"
	// global KV store - it should do for handlers
	usrs := make(map[uuid.UUID]storage.User)
	authID, _ := uuid.NewRandom()
	unauthID, _ := uuid.NewRandom()
	usrs[authID] = storage.User{Name: "frodo", APIkey: "1337"}
	m.users = usrs

	a := app{users: m}
	tCases := []struct {
		name     string
		expCode  int
		userID   uuid.UUID
		wantBody string
	}{
		{name: "succesful request - authenticated user",
			expCode:  200,
			userID:   authID,
			wantBody: `{"username":"frodo","APIkey":"1337"}`,
		},

		{name: "authed user - no info",
			expCode:  404,
			userID:   unauthID,
			wantBody: `{"error":"not found"}`},
		{
			name:     "no userID (or nil) in context would panic - should not occur",
			expCode:  401,
			userID:   uuid.UUID{},
			wantBody: `{"error":"you must be authenticated to access this resource"}`,
		},
	}
	for _, tc := range tCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, endpoint, nil)
			if err != nil {
				t.Fatal("can't make request")
			}
			ctx := context.WithValue(context.Background(), userIDctx, tc.userID)
			req = req.WithContext(ctx)
			rr := httptest.NewRecorder()
			a.GetUserData(rr, req)
			got := rr.Body.String()
			assertStatus(t, rr.Code, tc.expCode)
			if rr.Result().Header.Get("content-type") != "application/json" {
				t.Errorf("failed to set content header, got: %v", rr.Result().Header)
			}
			if got != tc.wantBody {
				t.Errorf("GET %v: want: %v, got: %v", endpoint, tc.wantBody, string(got))
			}
		})
	}
}

func TestCreateUser(t *testing.T) {
	endpoint := "/v1/users"
	m := &MockUserStore{num: 3}
	a := app{users: m}
	type resp struct {
		Username string `json:"name,omitempty"`
		APIkey   string `json:"APIkey,omitempty"`
		Error    string `json:"error,omitempty"`
	}
	type req struct {
		Username string
	}
	testCases := []struct {
		name     string
		expCode  int
		input    string
		wantBody string
	}{{
		name:     "empty json object",
		expCode:  http.StatusBadRequest,
		input:    `{}`,
		wantBody: `{"error":"body must not be empty"}`},
		{name: "invalid request - wrong json key",
			expCode:  http.StatusBadRequest,
			input:    `{"bla":"frodo"}`,
			wantBody: `{"error":"body has unsupported key \"bla\""}`,
		},
		{name: "boggus body",
			expCode:  http.StatusBadRequest,
			input:    `#{":f%#¢`,
			wantBody: `{"error":"poorly-formed JSON (at position 1)"}`, // flaky
		},
		{name: "well formed request",
			expCode:  http.StatusCreated,
			input:    `{"name":"frodo"}`,
			wantBody: `{"username":"frodo","APIkey":"1337"}`, // TODO establish name vs username
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b := strings.NewReader(tc.input) // TODO convert to json?
			req, err := http.NewRequest(http.MethodPost, endpoint, b)
			if err != nil {
				t.Fatal("can't make request")
			}
			rr := httptest.NewRecorder()
			a.CreateUser(rr, req)
			got := rr.Body.String()
			assertStatus(t, rr.Code, tc.expCode)
			if got != tc.wantBody {
				t.Errorf("POST %v: want: %v, got: %v", endpoint, tc.wantBody, string(got))
			}
			if rr.Result().Header.Get("content-type") != "application/json" {
				t.Errorf("failed to set content header, got: %v", rr.Result().Header)
			}

		})
	}
	t.Run("no body - should not panic", func(t *testing.T) {
		expCode := http.StatusBadRequest
		req, err := http.NewRequest(http.MethodPost, endpoint, nil)
		if err != nil {
			t.Fatal("can't make request")
		}
		rr := httptest.NewRecorder()

		func() {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("sending request panicked handler, but I expected no panic")
				}
			}()
			a.CreateUser(rr, req)
		}()
		want := `{"error":"body must not be empty"}`
		got := rr.Body.String()
		assertStatus(t, expCode, rr.Code)
		if got != want {
			t.Errorf("POST %v: want: %v, got: %v", endpoint, want, string(got))
		}
	})

}
func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf(" [ STATUS CODE MISMATCH ] - got: %v, wanted: %v", got, want)
	}
}
func (m MockUserStore) Create(ctx context.Context, username string) (storage.User, error) {
	return storage.User{Name: "frodo", APIkey: "1337"}, nil
}
func (m MockUserStore) WhoIs(context.Context, string) (storage.User, error) {
	return storage.User{}, nil
}
func (m MockUserStore) Get(ctx context.Context, userID uuid.UUID) (storage.User, error) {
	u, ok := m.users[userID]
	if !ok {
		return storage.User{}, storage.ErrNotFound
	}
	return u, nil
}
