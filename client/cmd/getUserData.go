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

// getUserDataCmd represents the getUserData command
var getUserDataCmd = &cobra.Command{
	Use:   "getUserData",
	Short: "retrieve user's feeds, username, apikey",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		apiKey, err := ReadApiKey(DEFAULT_ENV_FILE)
		if err != nil {
			return fmt.Errorf("cannot load apikey: %v", err)
		}
		//rooturl := API_URL
		return getUserDataAction(os.Stdout, API_URL, apiKey)
	},
}

func init() {
	rootCmd.AddCommand(getUserDataCmd)
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
