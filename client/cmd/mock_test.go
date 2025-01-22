package cmd

import (
	"net/http"
	"net/http/httptest"
)

const (
	createUserSuccess = iota
	getUsersData
	createFeedSuccess
	getFeeds
	getFeedsNoutFound
	followExistingFeed
	deleteFollowedFeed
	allPosts
	malformedRequest
	notFound
)

// all responses accounted for?
// testResp simulates test reponses from the API
var testResp = map[int]struct {
	Status int
	Body   string
}{
	createUserSuccess: {
		Status: http.StatusOK,
		Body:   `{"name":"Frodo","apikey":"bXkgcHJlY2lvdXM-aXRzLW1pbmU-bXkgZGVhciBnYW5kYWxm"}`},

	createFeedSuccess: {
		Status: http.StatusOK,
		Body:   `{"name":"Blog on basting","url":"www.kitchen-baste.com/xml","id":"f3ffd9ef-69bd-4f28-9cee-3c6cbbfacb3e"}`},

	getFeeds: {
		Status: http.StatusOK,
		Body: `[
	{"name":"Blog on basting","url":"www.kitchen-baste.com/xml","id":"f3ffd9ef-69bd-4f28-9cee-3c6cbbfacb3e"},{"name":"Blog on compression","url":"www.techsavy.com/xml","id":"1eb60252-3712-4263-bdca-de7a3b6825e2"}
	]`},

	getFeedsNoutFound: {
		Status: http.StatusNotFound,
		Body:   `{}`},

	followExistingFeed: {
		Status: http.StatusOK,
		Body:   `{"name":"Blog on compression","url":"www.techsavy.com/xml","id":"1eb60252-3712-4263-bdca-de7a3b6825e2"}`,
	},

	getUsersData: {
		Status: http.StatusOK,
		Body: `{"name":"Frodo","apikey":"bXkgcHJlY2lvdXM-aXRzLW1pbmU-bXkgZGVhciBnYW5kYWxm", "followedFeeds":[
	{"name":"Blog on basting","url":"www.kitchen-baste.com/xml","id":"f3ffd9ef-69bd-4f28-9cee-3c6cbbfacb3e"},{"name":"Blog on compression","url":"www.techsavy.com/xml","id":"1eb60252-3712-4263-bdca-de7a3b6825e2"}]}`,
	},

	deleteFollowedFeed: {
		Status: http.StatusOK,
		Body:   `{"Unfollowed feed"}`,
	},

	notFound: {
		Status: http.StatusNotFound,
		Body:   "404 - not found",
	},

	malformedRequest: {
		Status: http.StatusBadRequest,
		Body:   "",
	},

	allPosts: {
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
