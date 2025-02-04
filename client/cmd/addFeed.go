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
var addCmd = &cobra.Command{
	Use:     "add <feedName> <feedUrl>",
	Short:   "Add a feed from which you would like to collect posts",
	Example: "add 'memes' https://xkcd.com/rss.xml",
	Long: `Add a feed which you want to follow and receive posts/podcasts from.
	To provide a feed name using white space use single or double quotes
	
	Example:
	 add "funny memes" https://xkcd.com/rss.xml
	`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		apiKey, err := ReadAPIKey(credentialsFile)
		if err != nil {
			fmt.Fprintf(os.Stdout, "cannot read apiKey: %v", err)
			os.Exit(1)
		}
		return addFeedAction(os.Stdout, args, API_URL, apiKey)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
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
	if ok := isURL(feed); !ok {
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
	resp, err := sendReq(url, http.MethodPost, apiKey, "application/json", http.StatusCreated, &body)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
