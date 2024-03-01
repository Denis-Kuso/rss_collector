/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/Denis-Kuso/cli_rss/pkg/fetch"
)

var username string;
// createUserCmd represents the createUser command
var createUserCmd = &cobra.Command{
	Use:   "createUser <username>",
	Short: "Creates a new user.",
	Long: `Succesful creations returns a key, which is used
	for automatic login.`,
	Args: cobra.ExactArgs(1), //ARGS AND FLAGS ARE NOT THE SAME THING
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("createUser called with:", args)
		fetch.CreateUser(username)
	},
}

func init() {
	rootCmd.AddCommand(createUserCmd)
}
func CreateUser(username string) (apiKey string, err error) {
// request to server at createUser 
    ENDPOINT := "/users"
    data, err := fetchEndpoint(c, URL + ENDPOINT)
    if err != nil {
        fmt.Printf("ERR: %v\n", err)
	return "", err
    }
    // save api key/display apiKey 
    fmt.Printf("Got data: %v\n", string(data))
    return apiKey, nil
}
