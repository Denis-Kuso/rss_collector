/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/cobra"
)

// getFollowedFeedsCmd represents the getFollowedFeeds command
var getFollowedFeedsCmd = &cobra.Command{
	Use:   "getFollowedFeeds",
	Short: "Get all the feeds currently following.",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

	SOME example of usage: bla bla.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("getFollowedFeeds called")
	},
}

func init() {
	rootCmd.AddCommand(getFollowedFeedsCmd)
}

// TODO
func getAllFollowedFeedsAction(out io.Writer, rootURL, apiKey string) error {
	resp, err := getAllFollowedFeeds(rootURL, apiKey)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(out, string(resp))
	return err
}

func getAllFollowedFeeds(rootURL, apiKey string) ([]byte, error) {
	ENDPOINT := "/feed_follows"
	url := rootURL + ENDPOINT
	resp, err := sendReq(url, http.MethodGet, apiKey, "", http.StatusOK, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
