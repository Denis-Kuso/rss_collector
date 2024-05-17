/*
Copyright © 2024 Denis <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

// followFeedCmd represents the followFeed command
var followFeedCmd = &cobra.Command{
	Use:   "followFeed <feed_id>",
	Short: "Follow a feed added by someone else",
	Long: `There might be feeds on the server added by other users or offered
as potentialy interesting feeds. By browsing through them, any feed of interest
	can then be followed by using this command with a feed_id.`,
	Example: "followFeed c607531a-832a-4b44-9b13-3acd9839d165",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		apikey, err := ReadApiKey(credentialsFile)
		if err != nil {
			return fmt.Errorf("cannot read apiKey: %v\n", err)
		}
		return followFeedAction(os.Stdout, args[0], API_URL, apikey)
	},
}

func init() {
	rootCmd.AddCommand(followFeedCmd)

}

func followFeedAction(out io.Writer, feed_id, rootURL, apiKey string) error {
	resp, err := followFeed(feed_id, rootURL, apiKey)
	if err != nil {
		return err //to where?
	}
	_, err = fmt.Fprint(out, string(resp))
	return err
}

func followFeed(feed_id, url, apiKey string) ([]byte, error) {
	ENDPOINT := "/feed_follows"
	url += ENDPOINT
	if ok := isValidID(feed_id); !ok {
		return nil, fmt.Errorf("invalid id format: %v", feed_id)
	}
	feed := struct {
		FeedID string `json:"feed_id"`
	}{
		FeedID: feed_id}
	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(feed); err != nil {
		return nil, err
	}
	resp, err := sendReq(url, http.MethodPost, apiKey, "application/json", http.StatusOK, &body)
	if err != nil {
		// TODO add more error granularity
		fmt.Printf("ERR: %v\n", err)
		os.Exit(1)
	}
	return resp, nil
}
