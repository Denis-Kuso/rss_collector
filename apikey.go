package main

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/joho/godotenv"

)

// Extracts apiKey from response
// example {"name":"user", "apiKey": "1337"}
// return "1337"
// I could return []bytes as well
func ExtractApiKey(body []byte) (string, error) {
	type resp struct {
		ApiKey string `json:"ApiKey"`
	}
	var apiKey string
	r := resp{}
	err := json.Unmarshal(body, &r)
	if err != nil {
		return apiKey, fmt.Errorf("cannot extract ApiKey: %w", err) 
	}
	return r.ApiKey, nil
}

// reads apikey from filename
func ReadApiKey(filename string) (string, error) {
	const APIKEY string = "APIKEY"
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

// saves apiKey to disk
func SaveApiKey(apiKey []byte, out io.Writer) error {	
	n, err := out.Write(apiKey)
	if err != nil {
		return err
	}
	if n < len(apiKey){
		return fmt.Errorf("partial write: %d bytes written", n)
	}
	return nil
}
