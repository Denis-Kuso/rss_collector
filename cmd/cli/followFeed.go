package main

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
var followCmd = &cobra.Command{
	Use:   "follow <feedID>",
	Short: "Follow a feed added by someone else",
	Long: `There might be feeds on the server added by other users or offered
as potentialy interesting feeds. By browsing through them, any feed of interest
	can then be followed by using this command with a feedID.`,
	Example: "follow c607531a-832a-4b44-9b13-3acd9839d165",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		apikey, err := ReadAPIKey(credentialsFile)
		if err != nil {
			return fmt.Errorf("cannot read apiKey: %v", err)
		}
		return followFeedAction(os.Stdout, args[0], API_URL, apikey)
	},
}

func init() {
	rootCmd.AddCommand(followCmd)

}

func followFeedAction(out io.Writer, feedID, rootURL, apiKey string) error {
	resp, err := followFeed(feedID, rootURL, apiKey)
	if err != nil {
		return err
	}
	_, err = fmt.Fprint(out, string(resp))
	return err
}

func followFeed(feedID, url, apiKey string) ([]byte, error) {
	ENDPOINT := "/feed_follows"
	url += ENDPOINT
	if ok := isValidID(feedID); !ok {
		return nil, fmt.Errorf("invalid id format: %v", feedID)
	}
	feed := struct {
		FeedID string `json:"feedID"`
	}{
		FeedID: feedID}
	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(feed); err != nil {
		return nil, err
	}
	resp, err := sendReq(url, http.MethodPost, apiKey, "application/json", http.StatusOK, &body)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
