package cmd

import (
	"net/http"
	"net/http/httptest"
)

const (
	CREATE_USER_SUCCESS = iota
	CREATE_FEED_SUCCESS 
	GET_FEEDS 
	GET_FEEDS_NOT_FOUND 
	MALFORMERD_REQUEST 
	ROOT
	UNAUTHORISED
	NO_HEADER
	NOT_FOUND
	ALL_POSTS
	CREATED
)

// testResp simulates test reponses from the API
var testResp = map[int]struct {
	Status int
	Body   string
}{
	CREATE_USER_SUCCESS : {
		Status: http.StatusOK,
		Body: `{"ID":"001","CreatedAt":"testTime","UpdatedAt":"2019-10-28T08:23:38.310097076-04:00","Name":"TestName","ApiKey":"414141414141"}`},

	CREATE_FEED_SUCCESS: {
		Status: http.StatusOK,
		Body: `{
			"feed": {"id":"1", "CreatedAt":"someTime",
					"updatedAt":"someTime",
					"name":"testName",
					"url":"testingURL",
					"userID":"testID",
					"LastFetchedAt":"someTime"},
			"feedFollow: {"ID":"testId",
					"CreatedAt":"testTime",
					"UpdatedAt":"testTime",
					"UserID":"testID",
					"FeedID": "testID"}	
}`},
	GET_FEEDS: {
		Status: http.StatusOK,
		Body: `[
	":{
	"ID": "some_id",
	"CreatedAt": "some_time",
	"UpdatedAt": "Sine_time",
	"Name": "some_name",
	"Url": "someUrl",
	"UserID": "someid",
	"LastFetchedAt": "someTime"
	}`},

	GET_FEEDS_NOT_FOUND: {
		Status: http.StatusNotFound,
		Body: `{
}`},

	ROOT: {
		Status: http.StatusOK,
		Body:   "welcome Gandalf",
	},

	UNAUTHORISED: {
		Status: http.StatusUnauthorized,
		Body: `{"error":"Unauthorized"}`,
		},
	
	NO_HEADER: {
			Status: http.StatusBadRequest,
			Body: `{"error":"No header included"}`,
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
		Body: "",
		},
	ALL_POSTS: {
		Status: http.StatusOK,
		Body: `[{
				"ID":"some_id",
				"CreatedAt":"Sometime",
				"UpdatedAt":"SomeTime",
				"Title":"SomeTitle",
				"Url":"SomeUrl",
				"Description":"someString",
				"PublishedAt":"Sometime",
				"FeedID":"someID"
				}]`},
}

// mockServer creates a mock server to simulate the RSS API
func mockServer(h http.HandlerFunc) (string, func()) {
	ts := httptest.NewServer(h)

	return ts.URL, func() {
		ts.Close()
	}
}
