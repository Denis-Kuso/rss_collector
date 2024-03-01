package cmd

import (
	"bytes"
	//"errors"
	"fmt"

	"net/http"
	"testing"
)

func TestAddFeed(t *testing.T) {
	testCases := []struct {
		name     string
		expError error
		expOut   string
		feedName string
		feedURL  string
		resp     struct {
			Status int
			Body   string
		}
		closeServer bool
	}{
		{name: "add valid feed",
			expError: nil,
			expOut: `{
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
}`,
			feedName: "testName",
			feedURL:  "testingURL",
			resp:     testResp["New feed: valid req"],
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// todo validate request
			url, cleanup := mockServer(
				func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(tc.resp.Status)
					fmt.Fprintln(w, tc.resp.Body)
				})
			defer cleanup()
			fmt.Printf("Using temp url: %v\n", url)
			// test function
			var out bytes.Buffer
			if err := addFeedAction(&out, []string{tc.feedName, tc.feedURL}, url); err != nil {
				if tc.expError == nil {
					t.Fatalf("Expected no error, got %q.\n", err)
				}
				if tc.expError != err {
					t.Errorf("Expected: %v, got %v\n", tc.expError, err)
				}
			}
			if tc.expOut != out.String() {
				t.Errorf("Expected: %q, \n\tgot: %q\n", tc.expOut, out.String())
			}
		})
	}
}

//func TestGetFeeds(t *testing.T) {
//	testCases := []struct {
//		name     string
//		expError error
//		expOut   string
//		resp     struct {
//			Status int
//			Body   string
//		}
//		closeServer bool
//	}{
//		{name: "Get existing feeds",
//			expError: nil,
//			expOut:   "-  1  Task 1\n-  2  Task 2\n",
//			resp:     testResp["Get - feeds: exists"],
//		},
//		{name: "NoResults",
//			expError: ErrNotFound,
//			resp:     testResp["Get - feeds: not found"]},
//		{name: "InvalidURL",
//			expError:    ErrConnection,
//			resp:        testResp["notFound"],
//			closeServer: true},
//	}
//
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			url, cleanup := mockServer(
//				func(w http.ResponseWriter, r *http.Request) {
//					w.WriteHeader(tc.resp.Status)
//					fmt.Fprintln(w, tc.resp.Body)
//				})
//			defer cleanup()
//
//			if tc.closeServer {
//				cleanup()
//			}
//
//			var out bytes.Buffer
//
//			err := getFeeds(&out, url)
//
//			if tc.expError != nil {
//				if err == nil {
//					t.Fatalf("Expected error %q, got no error.", tc.expError)
//				}
//
//				if !errors.Is(err, tc.expError) {
//					t.Errorf("Expected error %q, got %q.", tc.expError, err)
//				}
//				return
//			}
//
//			if err != nil {
//				t.Fatalf("Expected no error, got %q.", err)
//			}
//
//			if tc.expOut != out.String() {
//				t.Errorf("Expected output %q, got %q", tc.expOut, out.String())
//			}
//		})
//	}
//}
