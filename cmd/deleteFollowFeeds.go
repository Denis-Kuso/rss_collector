/*
Copyright © 2024 Denis<EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

// deleteFollowFeedsCmd represents the deleteFollowFeeds command
var deleteFollowFeedsCmd = &cobra.Command{
	Use:   "unfollowFeed <feedID>",
	Short: "Stop following a feed",
	Long: `If you no longer wish to receive posts from a feed, use this command
	to stop following the feed. FeedID can be obtained by calling getFollowedFeeds,
	getUserData or getAllFeeds`,
	Example: "unfollowFeed c607531a-832a-4b44-9b13-3acd9839d165",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		apiKey, err := ReadApiKey(DEFAULT_ENV_FILE)
		if err != nil {
			return fmt.Errorf("cannot load apikey: %v", err)
		}
		return deleteFollowFeedAction(os.Stdout, ROOT_URL, apiKey, args[0])
	},
}

func init() {
	rootCmd.AddCommand(deleteFollowFeedsCmd)
}

func deleteFollowFeedAction(out io.Writer, rootURL, apiKey, feedFollowID string) error {
	resp, err := deleteFollowFeed(rootURL, apiKey, feedFollowID)
	if err != nil {
		// which errors are possible?
		return fmt.Errorf("cannot unfollow feed: \"%v\" due to: %v", feedFollowID, err)
	}
	_, err = fmt.Fprintf(out, string(resp))
	return err
}

func deleteFollowFeed(rootURL, apiKey, feedFollowID string) ([]byte, error) {
	ENDPOINT := "/feed_follows/"
	if ok := isValidID(feedFollowID); !ok {
		return nil, fmt.Errorf("invalid id format: %v", feedFollowID)
	}
	url := rootURL + ENDPOINT + feedFollowID
	resp, err := sendReq(url, http.MethodDelete, apiKey, "", http.StatusOK, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
