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
	Short: "Follow a feed added by someone else.",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
to quickly create a Cobra application.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		rootURL := ROOT_URL 
		apikey, err := ReadApiKey(DEFAULT_ENV_FILE)
		if err != nil {
			fmt.Fprintf(os.Stdout, "cannot read apikey: %v", err)
		}
		followFeedAction(os.Stdout, args[0], rootURL, apikey)
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
	if ok := validFeedID(feed_id); !ok {
		return nil, nil // TODO some error, should this check happen before?
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
	fmt.Printf("Successfuly followed feed: %s\n", feed)
	return resp, nil

}

// TODO implement
func validFeedID(feedID string) bool {
	return true
}
