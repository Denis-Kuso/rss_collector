/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// createUserCmd represents the createUser command
var createUserCmd = &cobra.Command{
	Use:   "createUser --username <username>",
	Short: "Creates a new user.",
	Long: `MUCH longer description`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("createUser called")
	},
}

func init() {
	rootCmd.AddCommand(createUserCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createUserCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createUserCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
