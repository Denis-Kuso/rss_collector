/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

// getFollowedFeedsCmd represents the getFollowedFeeds command
var getFollowedFeedsCmd = &cobra.Command{
	Use:   "getFollowedFeeds",
	Short: "Get all the feeds currently following.",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

	SOME example of usage: bla bla.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		apiKey, err := ReadApiKey(DEFAULT_ENV_FILE)
		if err != nil {
			fmt.Fprintf(os.Stdout, "cannot load apikey: %v", err)
			os.Exit(5)
		}
		return getAllFollowedFeedsAction(os.Stdout, ROOT_URL, apiKey)
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
