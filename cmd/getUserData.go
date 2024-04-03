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
	Short: "Output user data.",
	Long:  `Perhaps unnessacry to use long description.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("getUserData called")
		apiKey, err := ReadApiKey(DEFAULT_ENV_FILE)
		if err != nil {
			fmt.Fprintf(os.Stdout, "cannot load apikey: %v", err)
			os.Exit(5)
		}
		return getUserDataAction(os.Stdout, ROOT_URL, apiKey)
	},
}

func init() {
	rootCmd.AddCommand(getUserDataCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getUserDataCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getUserDataCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
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
