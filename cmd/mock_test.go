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
		"Create User: resultsMany": {
		Status: http.StatusOK,
		Body: `{
  "results": 
    {
      "ID": "Task 1",
      "CreatedAt": false,
      "UpdatedAt": "2019-10-28T08:23:38.310097076-04:00",
      "Name": "0001-01-01T00:00:00Z"
	}

}`},
	"Get - feeds: exists": {
		Status: http.StatusOK,
		Body: `{
	"results":{
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
		Body:   "There's an API here",
	},

	"notFound": {
		Status: http.StatusNotFound,
		Body:   "404 - not found",
	},

	"created": {
		Status: http.StatusCreated,
		Body:   "",
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
