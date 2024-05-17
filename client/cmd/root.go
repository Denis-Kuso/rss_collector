/*
Copyright © 2024 Denis Kusic<EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cli_rss",
	Short: "a CLI client for rss feeds",
	Long: `Add feeds which you would like to follow, follow feeds added by
	other users. Unfollow them for whatever reason. Collect posts from the
	followed feeds`,
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	cobra.OnInitialize(initConfig)
	msg := fmt.Sprintf("config file (default is %s)", DEFAULT_ENV_FILE)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", msg)

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

var (
	// server address
	API_URL         string
	cfgFile         string
	credentialsFile string
)

const (
	DEFAULT_ENV_FILE      string = "./.env"
	defaultCredentialsLoc string = "./.credentials"
)

func initConfig() {
	const credentialsKey string = "CRED_LOC"
	const urlKey string = "SERVER_URL"
	if cfgFile == "" {
		cfgFile = DEFAULT_ENV_FILE
	}
	err := godotenv.Load(cfgFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	API_URL = os.Getenv(urlKey)
	if API_URL == "" {
		fmt.Printf("No server address specified: \"%s\"\n", API_URL)
		os.Exit(1)
	}
	credentialsFile = os.Getenv(credentialsKey)
	if credentialsFile == "" {
		credentialsFile = defaultCredentialsLoc
	}
}
