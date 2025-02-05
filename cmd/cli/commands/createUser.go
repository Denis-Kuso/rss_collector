package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

// createCmd represents the createUser command
var createCmd = &cobra.Command{
	Use:   "create <username>",
	Short: "Create a new user",
	Long: `A user needs to be created in order to use this tool.
	Succesful creation returns a key, enabling automatic login and usage
	of other commands. Further invocation without -o flag will alert the
	user that the existing key will be overwritten. If this is the desired
	behaviour and the user wants to retain access to previous the user, the
	key should be saved in a safe place.`,
	Example: "  create Frodo\nTo replace user:\n  create -o Smeagol",
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
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().BoolVarP(&Overwrite, "overwrite", "o", false, "overwrite an existing apikey")
}
func createUserAction(out io.Writer, baseURL, name string, overwrite bool) error {
	// prevent accidental overwriting of user
	if !overwrite {
		key, _ := ReadAPIKey(credentialsFile) // don't care about err, only if key is present, err return an empty string anyway
		if key != "" {
			return fmt.Errorf("a user already exists. Use -o flag if you would like to overwrite current user")
		}
	}

	resp, err := createUser(name, baseURL)
	if err != nil {
		return err
	}
	apikey, err := ExtractAPIKey(resp)
	if err != nil {
		return err
	}
	err = SaveAPIKeyF([]byte(apikey), credentialsFile)
	return displayUser(out, resp)
}

func displayUser(out io.Writer, body []byte) error {
	_, err := fmt.Fprintf(out, string(body))
	return err
}

func createUser(username, url string) (user []byte, err error) {
	ENDPOINT := "/users"
	url += ENDPOINT

	if ok := validateUsername(username); !ok {
		return nil, fmt.Errorf("%v: desired username: %q too long", ErrInvalidRequest, username)
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
	resp, err := sendReq(url, http.MethodPost, apiKey, "application/json", http.StatusCreated, &body)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return resp, nil
}
