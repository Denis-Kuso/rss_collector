/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	//"github.com/Denis-Kuso/cli_rss/pkg/fetch"
	"github.com/spf13/cobra"
)

var username string

const (
	MAX_USERNAME_LENGTH = 35
)

// createUserCmd represents the createUser command
var createUserCmd = &cobra.Command{
	Use:   "createUser <username>",
	Short: "Create a new user",
	Long: `A user needs to be created in order to use this tool.
	Succesful creation returns a key, enabling automatic login and usage
	of other commands. Further invocation without -o flag will alert the
	user that the existing key will be overwritten. If this is the desired
	behaviour and the user wants to retain access to previous the user, the
	key should be saved in a safe place.`,
	Example: "  createUser Frodo\nTo replace user:\n  createUser -o Smeagol",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		overwrite, err := cmd.Flags().GetBool("overwrite")
		if err != nil {
			return err
		}
		return createUserAction(os.Stdout, API_URL, args[0], overwrite)
	},
}

func init() {
	var Overwrite bool
	rootCmd.AddCommand(createUserCmd)
	createUserCmd.Flags().BoolVarP(&Overwrite, "overwrite", "o", false, "overwrite an existing apikey")
}
func createUserAction(out io.Writer, base_url, name string, overwrite bool) error {
	_, err := ReadApiKey(DEFAULT_ENV_FILE)
	// user already exists
	if err == nil && !overwrite {
		return fmt.Errorf("a user already exists. Use -o flag if you would like to overwrite current user")
	}

	resp, err := createUser(name, base_url)
	if err != nil {
		return err
	}
	apikey, err := ExtractApiKey(resp)
	if err != nil {
		return err
	}
	// TODO propagate error for display?
	err = SaveApiKeyF([]byte(apikey), DEFAULT_ENV_FILE)
	return displayUser(out, resp)
}

func displayUser(out io.Writer, body []byte) error {
	// verbose option
	_, err := fmt.Fprintf(out, string(body))
	return err
}
func createUserF(username, url string) (resp []byte, err error) {
	ENDPOINT := "/users"
	url += ENDPOINT
	data, err := fetchEndpoint(c, url)
	if err != nil {
		return nil, fmt.Errorf("ERR: %v, during fetching with url:%v\n", err, url)
	}
	return data, nil
}

func createUser(username, url string) (user []byte, err error) {
	ENDPOINT := "/users"
	url += ENDPOINT

	if ok := validateUsername(username); !ok {
		return nil, fmt.Errorf("%v: desired username: \"%s\" too long", ErrInvalidRequest, username)
	}
	name := struct {
		Username string `json:"name"`
	}{
		Username: username}
	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(name); err != nil {
		return nil, err
	}
	apiKey := ""
	resp, err := sendReq(url, http.MethodPost, apiKey, "application/json", http.StatusOK, &body)
	if err != nil {
		fmt.Printf("ERR from sendReq: %v\n", err)
		os.Exit(1)
	}
	// TODO Then what?? process response or pass response forward?
	return resp, nil
}
