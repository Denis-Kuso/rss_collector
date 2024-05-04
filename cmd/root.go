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
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cli_rss.yaml)")

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// website address specified in env file
var API_URL string

func initConfig() {
	const keyURL string = "SERVER_URL"

	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	API_URL = os.Getenv(keyURL)
	if API_URL == "" {
		fmt.Printf("No url specified: \"%s\"\n", API_URL)
		os.Exit(1)
	}
}
