package cmd

import (
	"bytes"
	//"errors"
	"fmt"
	"io"
	"log"
	"strings"

	"net/http"
	"testing"
)

type ExpReq struct {
	ExpContentType string
	ExpBody        string
	ExpUrlPath     string
	ExpAuthMethod  string
	ExpHTTPMethod  string
}

type Resp struct {
	Status int
	Body   string
}
type TestCase struct {
	name     string
	expError error
	expOut   string
	expReq   ExpReq
	limit    string
	feedID   string
	feedName string
	feedURL  string
	username string
	resp     struct {
		Status int
		Body   string
	}
}

func checkReq(t *testing.T, e ExpReq, tc TestCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != (e.ExpUrlPath) {
			t.Errorf("Poor request! Expected path: %q, got: %q", e.ExpUrlPath, r.URL.Path)
		}
		if r.Method != e.ExpHTTPMethod {
			t.Errorf("Poor request! Expected method: %q, got: %q", e.ExpHTTPMethod, r.Method)
		}
		authValue := r.Header.Get("Authorization")
		if authValue == "" {
			t.Fatal("No header provided")
		}
		authMethod := strings.Split(authValue, " ")[0]
		if authMethod != e.ExpAuthMethod {
			t.Fatalf("Incorrect authorization method, expected: %v, got: %v\n", e.ExpAuthMethod, authMethod)
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		r.Body.Close()
		if string(body) != e.ExpBody {
			t.Errorf("Poor request! Expected body: %q, got: %q", e.ExpBody, string(body))
		}
		contentType := r.Header.Get("Content-Type")
		if contentType != e.ExpContentType {
			t.Errorf("Poor request! Expected Content-Type: %q, got: %q", e.ExpContentType, contentType)
		}
		w.WriteHeader(tc.resp.Status)
		fmt.Fprintln(w, tc.resp.Body)
	}
}
func TestGetUserData(t *testing.T) {
	e := ExpReq{
		ExpUrlPath:     "/users",
		ExpHTTPMethod:  http.MethodGet,
		ExpAuthMethod:  "ApiKey",
		ExpContentType: "",
		ExpBody:        "",
	}
	ApiKey := "someFancy4|>11<3j" //TODO decide whether to test this and HOW
	testCases := []TestCase{
		{
			name:     "get_user_data",
			expError: nil,
			expReq:   e,
			expOut:   `{"ID":"someID","CreatedAt":"someTime","UpdatedAt":"someTime","Name":"testName","ApiKey":"1337"}` + string('\n'),
			resp:     testResp[GET_USERS_DATA],
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url, cleanup := mockServer(checkReq(t, tc.expReq, tc))
			defer cleanup()
			var out bytes.Buffer
			log.Printf("calling action with url: %v\n", url)
			err := getUserDataAction(&out, url, ApiKey)
			if err != nil {
				if tc.expError == nil {
					t.Fatalf("Expected no error, got: %q.\n", err)
				}
				if tc.expError != err {
					t.Errorf("Expected err: %v, got: %v\n", tc.expError, err)
				}
			}
			fmt.Println(out.String())
			if tc.expOut != out.String() {
				t.Errorf("Expected: %q, \n\tgot: %q\n", tc.expOut, out.String())
			}
		})
	}
}

func TestGetPosts(t *testing.T) {
	e := ExpReq{
		ExpUrlPath:     "/posts",
		ExpHTTPMethod:  http.MethodGet,
		ExpAuthMethod:  "ApiKey", // TODO perhaps use const or enum
		ExpContentType: "",
		ExpBody:        "",
	}
	ApiKey := "someFancy4|>11<3j" //TODO decide whether to test this and HOW
	testCases := []TestCase{
		{
			name:     "get posts - no limit provided",
			limit:    "",
			expError: nil,
			expOut: `[{
				"ID":"some_id",
				"CreatedAt":"Sometime",
				"UpdatedAt":"SomeTime",
				"Title":"SomeTitle",
				"Url":"SomeUrl",
				"Description":"someString",
				"PublishedAt":"Sometime",
				"FeedID":"someID"
				}]` + string('\n'),
			resp: testResp[ALL_POSTS],
		},
		{
			name:     "get posts - limit provided",
			limit:    "5",
			expError: nil,
			expOut: `[{
				"ID":"some_id",
				"CreatedAt":"Sometime",
				"UpdatedAt":"SomeTime",
				"Title":"SomeTitle",
				"Url":"SomeUrl",
				"Description":"someString",
				"PublishedAt":"Sometime",
				"FeedID":"someID"
				}]` + string('\n'),
			resp: testResp[ALL_POSTS],
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url, cleanup := mockServer(checkReq(t, e, tc))
			defer cleanup()
			var out bytes.Buffer
			err := getPostsAction(&out, url, ApiKey, tc.limit)
			if err != nil {
				if tc.expError == nil {
					t.Fatalf("Expected no error, got: %q.\n", err)
				}
				if tc.expError != err {
					t.Errorf("Expected err: %v, got: %v\n", tc.expError, err)
				}
			}
			if tc.expOut != out.String() {
				t.Errorf("Expected: %q, \n\tgot: %q\n", tc.expOut, out.String())
			}

		})
	}
}

func TestDeleteFeedFollow(t *testing.T) {
	e := ExpReq{
		ExpUrlPath:     "/feed_follows/",
		ExpHTTPMethod:  http.MethodDelete,
		ExpAuthMethod:  "ApiKey", // TODO perhaps use const or enum
		ExpContentType: "",
		ExpBody:        "",
	}
	ApiKey := "someFancy4|>11<3j" //TODO decide whether to test this and HOW
	testCases := []TestCase{
		{
			name:     "delete existing feed_follow",
			feedID:   "1337",
			expError: nil,
			expOut:   `{"Unfollowed feed"}` + string('\n'),
			resp:     testResp[DELETE_FOLLOW_FEED],
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url, cleanup := mockServer(checkReq(t, e, tc))
			defer cleanup()
			fmt.Print("Testing")
			fmt.Println(url)
			var out bytes.Buffer
			err := deleteFollowFeedAction(&out, url, ApiKey, tc.feedID)
			if err != nil {
				if tc.expError == nil {
					t.Fatalf("Expected no error, got: %q.\n", err)
				}
				if tc.expError != err {
					t.Errorf("Expected err: %v, got: %v\n", tc.expError, err)
				}
			}
			if tc.expOut != out.String() {
				t.Errorf("Expected: %q, \n\tgot: %q\n", tc.expOut, out.String())
			}
		})
	}
}

func TestGetAllFollowedFeeds(t *testing.T) {
	e := ExpReq{
		ExpUrlPath:     "/feed_follows",
		ExpHTTPMethod:  http.MethodGet,
		ExpAuthMethod:  "ApiKey", // TODO perhaps use const or enum
		ExpBody:        "",
		ExpContentType: "", // TODO Are these required
	}
	ApiKey := "someFancy4|>11<3j" //TODO decide whether to test this and HOW
	testCases := []TestCase{
		{
			name:     "Get all followed feeds: valid",
			expError: nil,
			expOut: `[
	{
	"ID": "some_id",
	"CreatedAt": "some_time",
	"UpdatedAt": "some_time",
	"Name": "some_name",
	"Url": "someUrl",
	"UserID": "someid",
	"LastFetchedAt": "someTime"
	}]` + string('\n'),
			resp: testResp[GET_FEEDS],
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// validate request well formed
			url, cleanup := mockServer(checkReq(t, e, tc))
			defer cleanup()
			var out bytes.Buffer
			err := getAllFollowedFeedsAction(&out, url, ApiKey)
			if err != nil {
				if tc.expError == nil {
					t.Fatalf("Expected no error, got: %q.\n", err)
				}
				if tc.expError != err {
					t.Errorf("Expected err: %v, got: %v\n", tc.expError, err)
				}
			}
			if tc.expOut != out.String() {
				t.Errorf("Expected: %q, \n\tgot: %q\n", tc.expOut, out.String())
			}
		})
	}
}
func TestFollowFeed(t *testing.T) {
	e := ExpReq{
		ExpUrlPath:     "/feed_follows",
		ExpHTTPMethod:  http.MethodPost,
		ExpAuthMethod:  "ApiKey", // TODO perhaps use const or enum,
		ExpContentType: "application/json",
		ExpBody:        `{"feed_id":"c5c9212c-57a3-4d68-b42e-addd951502c0"}` + string('\n'), //Encoder adds a newline char
	}
	ApiKey := "someFancy4|>11<3j" //TODO decide whether to test this and HOW
	testCases := []TestCase{
		{
			name:     "Follow existing feed",
			expError: nil,
			expOut:   `{"ID":"c52d3a13-2245-4991-8012-8856417b706f","CreatedAt":"2024-02-26T17:47:09.099267Z","UpdatedAt":"2024-02-26T17:47:09.099268Z","UserID":"8f588151-5489-4668-bfff-8c50021c1160","FeedID":"c5c9212c-57a3-4d68-b42e-addd951502c0"}` + string('\n'),
			feedID:   "c5c9212c-57a3-4d68-b42e-addd951502c0",
			resp:     testResp[FOLLOW_EXISTING_FEED]},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// validate request well formed
			url, cleanup := mockServer(checkReq(t, e, tc))
			defer cleanup()
			var out bytes.Buffer
			err := followFeedAction(&out, tc.feedID, url, ApiKey)
			if err != nil {
				if tc.expError == nil {
					t.Fatalf("Expected no error, got: %q.\n", err)
				}
				if tc.expError != err {
					t.Errorf("Expected err: %v, got: %v\n", tc.expError, err)
				}
			}
			if tc.expOut != out.String() {
				t.Errorf("Expected: %q, \n\tgot: %q\n", tc.expOut, out.String())
			}
		})
	}
}

func TestCreateUser(t *testing.T) {
	e := ExpReq{
		ExpUrlPath:     "/users",
		ExpHTTPMethod:  http.MethodPost,
		ExpContentType: "application/json",
		ExpBody:        `{"name":"testUsername"}` + string('\n'), //Encoder adds a newline char
	}
	testCases := []TestCase{
		{
			name:     "Valid request",
			expError: nil,
			expOut: `{"ID":"001","CreatedAt":"testTime","UpdatedAt":"2019-10-28T08:23:38.310097076-04:00","Name":"TestName","ApiKey":"414141414141"}
`,
			username: "testUsername",
			resp:     testResp[CREATE_USER_SUCCESS],
		},
		//			{testName: "Invalid request",
		//			expError: ErrInvalidRequest,
		//			expOut: "\n",
		//			username: "testUsername",
		//			resp: testResp[MALFORMERD_REQUEST],
		//	}	,
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// validate request well formed
			url, cleanup := mockServer(checkReq(t, e, tc))
			defer cleanup()
			// tests
			var out bytes.Buffer
			err := createUserAction(&out, url, tc.username, true)
			if err != nil {
				if tc.expError == nil {
					t.Fatalf("Expected no error, got: %q.\n", err)
				}
				if tc.expError != err {
					t.Errorf("Expected err: %v, got: %v\n", tc.expError, err)
				}
			}
			if tc.expOut != out.String() {
				t.Errorf("Expected: %q, \n\tgot: %q\n", tc.expOut, out.String())
			}
		})
	}
}

func TestAddFeed(t *testing.T) {
	e := ExpReq{
		ExpUrlPath:     "/feeds",
		ExpHTTPMethod:  http.MethodPost,
		ExpContentType: "application/json",
		ExpAuthMethod:  "ApiKey",
		ExpBody:        `{"name":"testName","url":"testingURL"}` + string('\n'), //Encoder adds a newline char
	}
	// TODO check header is present?AUTH...
	testCases := []TestCase{
		{
			name:     "add valid feed",
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
}` + string('\n'), //Encoder adds a newline char
			feedName: "testName",
			feedURL:  "testingURL",
			resp:     testResp[CREATE_FEED_SUCCESS],
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// validate request
			url, cleanup := mockServer(checkReq(t, e, tc))
			defer cleanup()
			// test function
			var out bytes.Buffer
			if err := addFeedAction(&out, []string{tc.feedName, tc.feedURL}, url, "1337"); err != nil {
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
