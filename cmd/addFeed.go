/*
Copyright © 2024 Denis Kusic<EMAIL ADDRESS>
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

// addFeedCmd represents the addFeed command
var addFeedCmd = &cobra.Command{
	Use:   "addFeed <feedName> <feedUrl>",
	Short: "Add a feed to the feeder",
	Long: `The feed added is automatically followed by the user.

 Example: 
	addFeed blogOnAgi www.agiblog.com`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		rooturl := ROOT_URL
		apiKey, err := ReadApiKey(DEFAULT_ENV_FILE)
		if err != nil {
			fmt.Fprintf(os.Stdout, "cannot read apiKey: %v", err)
			os.Exit(5)
		}
		return addFeedAction(os.Stdout, args, rooturl, apiKey)
	},
}

func init() {
	rootCmd.AddCommand(addFeedCmd)
}

func addFeedAction(out io.Writer, args []string, rooturl, apiKey string) error {
	name, feed := args[0], args[1]
	resp, err := addFeed(name, feed, rooturl, apiKey)
	if err != nil {
		return err
	}
	return displayAddFeed(out, resp)
}

// custom printing
func displayAddFeed(out io.Writer, feed []byte) error {
	_, err := fmt.Fprint(out, string(feed))
	return err
}

func addFeed(name, feed, url, apiKey string) ([]byte, error) {
	ENDPOINT := "/feeds"
	url += ENDPOINT
	if ok := isUrl(feed); !ok {
		return nil, fmt.Errorf("invalid url provided: %v", feed)
	}
	feedex := struct {
		Name string `json:"name"`
		Feed string `json:"url"`
	}{
		Name: name,
		Feed: feed}

	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(feedex); err != nil {
		return nil, err
	}
	resp, err := sendReq(url, http.MethodPost, apiKey, "application/json", http.StatusOK, &body)
	if err != nil {
		// TODO add more error granularity
		return nil, err
	}
	return resp, nil
}
