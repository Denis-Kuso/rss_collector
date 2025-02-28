package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/joho/godotenv"
)

const APIKEY string = "APIKEY"

func ExtractAPIKey(body []byte) (string, error) {
	type resp struct {
		APIKey string `json:"ApiKey"`
	}
	var apiKey string
	r := resp{}
	err := json.Unmarshal(body, &r)
	if err != nil {
		return apiKey, fmt.Errorf("cannot extract APIKey: %w", err)
	}
	return r.APIKey, nil
}

// reads apikey from filename
func ReadAPIKey(filename string) (string, error) {
	envKeys, err := godotenv.Read(filename)
	if err != nil {
		return "", fmt.Errorf("failed loading env file: %s, %w", filename, err)
	}
	apikey, ok := envKeys[APIKEY]
	if !ok {
		return "", fmt.Errorf("filename: %s contains no apikey", filename)
	}
	return apikey, nil
}

func SaveAPIKeyF(apiKey []byte, destName string) error {
	data := []byte(fmt.Sprintf("%s=%s", APIKEY, apiKey))
	err := os.WriteFile(destName, data, 0666)
	return err
}

// saves apiKey to disk
func SaveAPIKey(apiKey []byte, out io.Writer) error {
	prefix := []byte(APIKEY + "=")
	prefix = append(prefix, apiKey...)
	n, err := out.Write(prefix)
	if err != nil {
		return err
	}
	if n < len(apiKey) {
		return fmt.Errorf("partial write: %d bytes written", n)
	}
	return nil
}
