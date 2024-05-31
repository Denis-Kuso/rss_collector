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

// showCmd represents the getPublicFeeds command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show feeds already available to follow",
	Long: `If unsure what you might follow or how this tool works, using this command you
	will get some tech-related feeds which you can follow or use to try the tool.
	Does not require being a registered user`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return getPublicFeedsAction(os.Stdout, API_URL)
	},
}

func init() {
	rootCmd.AddCommand(showCmd)

}
func getPublicFeedsAction(out io.Writer, rootURL string) error {
	resp, err := getPublicFeeds(rootURL)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(out, string(resp))
	return err
}
func getPublicFeeds(rootURL string) ([]byte, error) {
	ENDPOINT := "/feeds"
	url := rootURL + ENDPOINT
	apiKey := ""
	resp, err := sendReq(url, http.MethodGet, apiKey, "", http.StatusOK, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
