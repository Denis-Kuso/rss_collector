//go:build integration

package cmd

import (
	"fmt"
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

func TestIntegration(t *testing.T) {
	fmt.Println("TESTING with real server...")
}
