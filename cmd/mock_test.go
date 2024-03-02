package cmd

import (
	"net/http"
	"net/http/httptest"
)

// testResp simulates test reponses from the API
var testResp = map[string]struct {
	Status int
	Body   string
}{
	"Create User: success": {
		Status: http.StatusOK,
		Body: `{"ID": "001","CreatedAt":"testTime","UpdatedAt": "2019-10-28T08:23:38.310097076-04:00","Name":"TestName","ApiKey":"414141414141"}`},

	"New feed: valid req": {
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
	"Get feeds": {
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

	"Get - feeds: not found": {
		Status: http.StatusNotFound,
		Body: `{
}`},

	"root": {
		Status: http.StatusOK,
		Body:   "welcome Gandalf",
	},

	"notFound": {
		Status: http.StatusNotFound,
		Body:   "404 - not found",
	},

	"created": {
		Status: http.StatusCreated,
		Body:   "",
	},
	"malformed request": {
		Status: http.StatusBadRequest,
		Body: "",
		},
	"retrieve posts": {
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
