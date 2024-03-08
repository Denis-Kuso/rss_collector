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

// getPostsCmd represents the getPosts command
var getPostsCmd = &cobra.Command{
	Use:   "getPosts",
	Short: "Get posts",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("getPosts called")
	},
}

func init() {
	rootCmd.AddCommand(getPostsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getPostsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getPostsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
func getPostsAction(out io.Writer, rootURL, apiKey, limit string) error {
	// validate limit (before making a request or call to getPosts
	ok := validLimit(limit)
	if !ok {
		fmt.Printf("ERR: Invalid query value:%v\n", limit)
		os.Exit(1)
	}
	resp, err := getPosts(rootURL, apiKey, limit)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(out, string(resp))
	return err
}

func getPosts(rootURL, apiKey, limit string) ([]byte, error) {
	ENDPOINT := "/posts"
	//assuming valid limit (if any) provided
	url := rootURL + ENDPOINT
	if limit != "" {
		url += "?limit=" + limit
	}
	resp, err := sendReq(url, http.MethodGet, apiKey, "", http.StatusOK, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// TODO implement
func validLimit(limit string) bool {
	return true
}
