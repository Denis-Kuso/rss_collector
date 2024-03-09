/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"net/http"

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
	Short: "Creates a new user.",
	Long: `Succesful creations returns a key, which is used
	for automatic login.`,
	Args: cobra.ExactArgs(1), //ARGS AND FLAGS ARE NOT THE SAME THING
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("createUser called with:", args)
		return createUserAction(os.Stdout, ROOT_URL, args[0])
	},
}

func init() {
	rootCmd.AddCommand(createUserCmd)
}
func createUserAction(out io.Writer, base_url, name string) error {
	resp, err := createUser(name, base_url)
	if err != nil {
		fmt.Printf("Failed creating user: %s.Err: %v\n", name, err)
		os.Exit(1)
	}
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
		return nil, fmt.Errorf("ERR: %v, during fetching with url:%v\n",err, url) 
	}
	return data, nil
}

func createUser(username, url string) (user []byte, err error) {
	ENDPOINT := "/users"
	url += ENDPOINT

	if ok := validateUsername(username); !ok{
		fmt.Println("NOT OK")
		return nil, ErrInvalidRequest// TODO CHANGE ERR VALUE
	}
	name := struct {
		Username string `json:"name"`
		}{
			Username:username}
	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(name); err != nil {
		return nil, err
	}
	apiKey := ""
	resp, err := sendReq(url, http.MethodPost, apiKey, "application/json", http.StatusOK, &body)
	if err != nil {
		fmt.Printf("ERR from sendReq: %v\n",err)
		os.Exit(1)
	}
	// TODO Then what?? process response or pass response forward?
	return resp, nil
}

// Checks the validity of the username provided.
// Max length of username is 35 characters. White space cannot be used.
func validateUsername(username string) bool {
	runeUsername := []rune(username)
	if len(runeUsername) > MAX_USERNAME_LENGTH {
		return false
	}
	return true
}
