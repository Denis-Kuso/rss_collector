package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

// fetchCmd represents the getPosts command
var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Retrieve posts from followed feeds",
	Long: `Retrieve posts from followed feeds. If no feeds followed or no posts
	are found an empty list is returned`,
	RunE: func(cmd *cobra.Command, args []string) error {
		apiKey, err := ReadAPIKey(credentialsFile)
		if err != nil {
			return fmt.Errorf("cannot load apikey: %v", err)
		}
		return getPostsAction(os.Stdout, API_URL, apiKey, "")
	},
}

func init() {
	rootCmd.AddCommand(fetchCmd)

}
func getPostsAction(out io.Writer, rootURL, apiKey, limit string) error {
	resp, err := getPosts(rootURL, apiKey, limit)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(out, string(resp))
	return err
}

func getPosts(rootURL, apiKey, limit string) ([]byte, error) {
	ENDPOINT := "/posts"
	url := rootURL + ENDPOINT
	if limit != "" {
		if ok := validLimit(limit); !ok {
			return nil, fmt.Errorf("%s: invalid limit: %s", ErrInvalidRequest, limit)
		}
		url += "?limit=" + limit
	}
	resp, err := sendReq(url, http.MethodGet, apiKey, "", http.StatusOK, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
