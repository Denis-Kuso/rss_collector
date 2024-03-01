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
		fmt.Println("addFeed called")
		rooturl := "NON-existing-url"
		return addFeedAction(os.Stdout, args, rooturl)
	},
}

func init() {
	rootCmd.AddCommand(addFeedCmd)
}

func addFeedAction(out io.Writer, args []string, rooturl string) error {
	name, feed := args[0], args[1]
	resp, err := addFeed(name, feed, rooturl)
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

func addFeed(name, feed, url string) ([]byte, error) {
	ENDPOINT := "/feeds"
	url += ENDPOINT
	// validate arg given
	cleanFeed := validateArg(feed)
	// TODO stopping point here if feed is invalid
	// TODO need to provide apiKey
	feedex := struct {
		Name string `json:"name"`
		Feed string `json:"url"`
	}{
		Name: name,
		Feed: cleanFeed}

	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(feedex); err != nil {
		return nil, err
	}
	resp, err := sendReq(url, http.MethodPost, "application/json", http.StatusOK, &body)
	if err != nil {
		// TODO add more error granularity
		fmt.Printf("ERR: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Successfuly added feed: %s\n", feed)
	return resp, nil
}

// TODO
func validateArg(arg string) string {
	return arg
}
