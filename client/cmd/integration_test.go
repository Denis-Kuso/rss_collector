package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"testing"
)

// flow of tests:
// createUser
// get user data?
// getFeeds (already available from server)
// follow an existing feed
// get user data
// add feed
// get user data
// get feedFollows
// delete currently followed feed
// get user data
// get posts

// Do it without creating a user first, need to implement
// saving apikey, then retrieving...
func TestIntegration(t *testing.T) {
	fmt.Println("TESTING with real server...")
	fmt.Println("Ou jea, let's expose some flaws")
	type User struct {
		name string
		key  string
	}
	user := User{
		name: "Pobro",
		key:  "6711c5359a5bb4a60bfd37113689bc003e128764d2599a7974fbc77e1580c27c",
	}
	t.Run("Get existing user data", func(t *testing.T) {
		expOut := `{"ID":"a42113b1-ed14-458a-aecb-c42c0c17fbda","CreatedAt":"2024-02-02T17:09:44.334957Z","UpdatedAt":":"2024-02-02T17:09:44.334957Z","Name":"Pobro","ApiKey":"6711c5359a5bb4a60bfd37113689bc003e128764d2599a7974fbc77e1580c27c"}`
		var out bytes.Buffer
		err := getUserDataAction(&out, API_URL, user.key)
		if err != nil {
			t.Fatalf("Expected no error, got: %v\n", err)
		}
		if expOut != out.String() {
			t.Fatalf("Expected output: %v, got: %s\n", expOut, out.String())
		}
	})
	t.Run("Get user feeds", func(t *testing.T) {
		expOut := `TODO`
		var out bytes.Buffer
		key := "6711c5359a5bb4a60bfd37113689bc003e128764d2599a7974fbc77e1580c27c"
		err := getUserDataAction(&out, API_URL, key)
		if err != nil {
			if errors.Is(err, ErrInvalidResponse) {
				t.Fatalf("Expected error: %v, got: %v\n", ErrInvalidResponse, err)
			}
		}
		if expOut != out.String() {
			t.Fatalf("Expected output: %v, got: %s\n", expOut, out.String())
		}
	})
	t.Run("Get feed_follows", func(t *testing.T) {
		expOut := `TODO`
		var out bytes.Buffer
		key := "6711c5359a5bb4a60bfd37113689bc003e128764d2599a7974fbc77e1580c27c"
		err := getAllFollowedFeedsAction(&out, API_URL, key)
		if err != nil {
			t.Fatalf("Expected error: %v, got: %v\n", ErrInvalidResponse, err)
		}
		if expOut != out.String() {
			t.Fatalf("Expected output: %v, got: %s\n", expOut, out.String())
		}
	})
	t.Run("Get posts - no limit provided", func(t *testing.T) {
		expOut := `TODO`
		var out bytes.Buffer
		key := "6711c5359a5bb4a60bfd37113689bc003e128764d2599a7974fbc77e1580c27c"
		err := getPostsAction(&out, API_URL, key, "")
		if err != nil {
			t.Fatalf("Expected error: %v, got: %v\n", ErrInvalidResponse, err)
		}
		if expOut != out.String() {
			t.Fatalf("Expected output: %v, got: %s\n", expOut, out.String())
		}
	})
	t.Run("Create feed", func(t *testing.T) {
		expOut := `TODO`
		var out bytes.Buffer
		args := []string{"feedname", "feedURL"} // TODO
		key := "6711c5359a5bb4a60bfd37113689bc003e128764d2599a7974fbc77e1580c27c"
		err := addFeedAction(&out, args, API_URL, key)
		if err != nil {
			t.Fatalf("Expected error: %v, got: %v\n", ErrInvalidResponse, err)
		}
		if expOut != out.String() {
			t.Fatalf("Expected output: %v, got: %s\n", expOut, out.String())
		}
	})
	// Get user data here?
	t.Run("Follow feed", func(t *testing.T) {
		expOut := `TODO`
		var out bytes.Buffer
		feed_id := "feedID"
		key := "6711c5359a5bb4a60bfd37113689bc003e128764d2599a7974fbc77e1580c27c"
		err := followFeedAction(&out, feed_id, API_URL, key)
		if err != nil {
			t.Fatalf("Expected error: %v, got: %v\n", ErrInvalidResponse, err)
		}
		if expOut != out.String() {
			t.Fatalf("Expected output: %v, got: %s\n", expOut, out.String())
		}
	})
	t.Run("Unfollow feed", func(t *testing.T) {
		expOut := `TODO`
		var out bytes.Buffer
		feed_id := "feedID"
		key := "6711c5359a5bb4a60bfd37113689bc003e128764d2599a7974fbc77e1580c27c"
		err := deleteFollowFeedAction(&out, API_URL, key, feed_id)
		if err != nil {
			t.Fatalf("Expected error: %v, got: %v\n", ErrInvalidResponse, err)
		}
		if expOut != out.String() {
			t.Fatalf("Expected output: %v, got: %s\n", expOut, out.String())
		}
	})
}
