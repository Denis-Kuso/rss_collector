package command

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

// infoCmd represents the getUserData command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Retrieve user's feeds, username, apikey",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		apiKey, err := ReadAPIKey(credentialsFile)
		if err != nil {
			return fmt.Errorf("cannot load apikey: %v", err)
		}
		return getUserDataAction(os.Stdout, API_URL, apiKey)
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}

func getUserDataAction(out io.Writer, rootURL, apiKey string) error {
	data, err := getUserData(rootURL, apiKey)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(out, string(data))
	return err
}

func getUserData(rootURL, apiKey string) ([]byte, error) {
	ENDPOINT := "/users"
	url := rootURL + ENDPOINT
	resp, err := sendReq(url, http.MethodGet, apiKey, "", http.StatusOK, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
