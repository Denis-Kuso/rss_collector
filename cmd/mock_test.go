package cmd

import (
	"net/http"
	"net/http/httptest"
)

const (
	CREATE_USER_SUCCESS = iota
	GET_USERS_DATA
	CREATE_FEED_SUCCESS
	GET_FEEDS
	GET_FEEDS_NOT_FOUND
	FOLLOW_EXISTING_FEED
	DELETE_FOLLOW_FEED
	ALL_POSTS
	MALFORMERD_REQUEST
	ROOT
	UNAUTHORISED
	NO_HEADER
	NOT_FOUND
	CREATED
)

// all responses accounted for?
// testResp simulates test reponses from the API
var testResp = map[int]struct {
	Status int
	Body   string
}{
	CREATE_USER_SUCCESS: {
		Status: http.StatusOK,
		Body: `{"name":"Frodo","apikey":"bXkgcHJlY2lvdXM-aXRzLW1pbmU-bXkgZGVhciBnYW5kYWxm"}`},

	CREATE_FEED_SUCCESS: {
		Status: http.StatusOK,
		Body: `{"name":"Blog on basting","url":"www.kitchen-baste.com/xml","id":"f3ffd9ef-69bd-4f28-9cee-3c6cbbfacb3e"}`},
	GET_FEEDS: {
		Status: http.StatusOK,
		Body: `[
	{"name":"Blog on basting","url":"www.kitchen-baste.com/xml","id":"f3ffd9ef-69bd-4f28-9cee-3c6cbbfacb3e"},{"name":"Blog on compression","url":"www.techsavy.com/xml","id":"1eb60252-3712-4263-bdca-de7a3b6825e2"}
	]`},

	GET_FEEDS_NOT_FOUND: {
		Status: http.StatusNotFound,
		Body: `{}`},

	FOLLOW_EXISTING_FEED: {
		Status: http.StatusOK,
		Body: `{"name":"Blog on compression","url":"www.techsavy.com/xml","id":"1eb60252-3712-4263-bdca-de7a3b6825e2"}`,
	},
	GET_USERS_DATA: {
		Status: http.StatusOK,
		Body:   `{"name":"Frodo","apikey":"bXkgcHJlY2lvdXM-aXRzLW1pbmU-bXkgZGVhciBnYW5kYWxm", "followedFeeds":[
	{"name":"Blog on basting","url":"www.kitchen-baste.com/xml","id":"f3ffd9ef-69bd-4f28-9cee-3c6cbbfacb3e"},{"name":"Blog on compression","url":"www.techsavy.com/xml","id":"1eb60252-3712-4263-bdca-de7a3b6825e2"}]}`,
	},
	DELETE_FOLLOW_FEED: {
		Status: http.StatusOK,
		Body:   `{"Unfollowed feed"}`,
	},

	ROOT: {
		Status: http.StatusOK,
		Body:   "welcome Gandalf",
	},

	UNAUTHORISED: {
		Status: http.StatusUnauthorized,
		Body:   `{"error":"Unauthorized"}`,
	},

	NO_HEADER: {
		Status: http.StatusBadRequest,
		Body:   `{"error":"No header included"}`,
	},

	NOT_FOUND: {
		Status: http.StatusNotFound,
		Body:   "404 - not found",
	},

	CREATED: {
		Status: http.StatusCreated,
		Body:   "",
	},
	MALFORMERD_REQUEST: {
		Status: http.StatusBadRequest,
		Body:   "",
	},
	ALL_POSTS: {
		Status: http.StatusOK,
		Body: `[{"feedName":"XKCD",
			"title": "Research Account",
			"url": "https://xkcd.com/2894/"},
			{"feedName":"newswire",
			"title":"reddit kills api support",
			"url":"https://netnewswire.blog/feed.xml"}]`},
}

// mockServer creates a mock server to simulate the RSS API
func mockServer(h http.HandlerFunc) (string, func()) {
	ts := httptest.NewServer(h)

	return ts.URL, func() {
		ts.Close()
	}
}
