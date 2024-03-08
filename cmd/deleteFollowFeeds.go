/*
Copyright © 2024 Denis<EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/cobra"
)

// deleteFollowFeedsCmd represents the deleteFollowFeeds command
var deleteFollowFeedsCmd = &cobra.Command{
	Use:   "deleteFollowFeeds <id>",
	Short: "Stop following feed",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Some example of usage.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("deleteFollowFeeds called")
	},
}

func init() {
	rootCmd.AddCommand(deleteFollowFeedsCmd)

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteFollowFeedsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteFollowFeedsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func deleteFollowFeedAction(out io.Writer, rootURL, apiKey, feedFollowID string) error {
	resp, err := deleteFollowFeed(rootURL, apiKey, feedFollowID)
	if err != nil {
	}
	_, err = fmt.Fprintf(out, string(resp))
	return err
}

func deleteFollowFeed(rootURL, apiKey, feedFollowID string) ([]byte, error) {
	ENDPOINT := "/feed_follows/"
	// what about id
	// assuming valid feedFollowID (integer)// TODO ENFORCE validity of ID
	url := rootURL + ENDPOINT + feedFollowID
	resp, err := sendReq(url, http.MethodDelete, apiKey, "", http.StatusOK, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
