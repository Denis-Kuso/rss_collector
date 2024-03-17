package main

import (
	"bytes"
	//"os"
	"testing"
)

func TestReadApiKey(t *testing.T) {
	testCases := []struct {
		testName string
		expErr   bool
		expOut   string
		input    string
	}{
		{
			testName: "both file and key present",
			expErr:   false,
			expOut:   "1337",
			input:    ".testenv.txt",
		},
		{testName: "file present, no apikey present",
			expErr: true,
			input:  ".testenv-no-key.txt",
		},
		{testName: "provided file does not exist",
			expErr: true,
			expOut: "",
			input:  "madeupfile.txt",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			var out string
			out, err := ReadApiKey(tc.input)
			if err != nil && !tc.expErr {
				t.Fatalf("Expected no err, got: %v\n", err)
			}
			if err == nil && tc.expErr {
				t.Fatalf("filename: %s, found api:%s, expected err, got: %v\n", tc.input, out, err)
			}
			if out != tc.expOut {
				t.Logf("Expected key: %v, got: %v", tc.expOut, out)
				t.Fail()
			}
		})
	}

}
func TestSaveApiKey(t *testing.T) {
	testCases := []struct {
		testName string
		expErr   bool
		expOut   string
		input    []byte
	}{
		{testName: "happy case",
			expErr: false,
			expOut: "1337",
			input:  []byte("1337"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			var out bytes.Buffer
			err := SaveApiKey(tc.input, &out)
			if err != nil {
				if !tc.expErr {
					t.Fatalf("Expected no err: %v, got: %v", tc.expErr, err) // TODO Could/Should define error type
				}
			}
			if out.String() != tc.expOut {
				t.Logf("Expected key: %v, got: %v", tc.expOut, out.String())
				t.Fail()
			}
		})
	}
}

func TestExtractApiKey(t *testing.T) {
	testCases := []struct {
		testName string
		expErr   bool
		expOut   string
		input    []byte
	}{
		{testName: "happy case",
			expErr: false,
			expOut: "1337",
			input:  []byte(`{"name":"username","apiKey":"1337"}`),
		},
		{testName: "invalid json",
			expErr: true,
			expOut: "",
			input:  []byte("Hello"),
		},
		{testName: "empty string",
			expErr: true,
			expOut: "",
			input:  []byte(""),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			out, err := ExtractApiKey(tc.input)
			if err != nil {
				if !tc.expErr {
					t.Fatalf("Expected no err: %v, got: %v", tc.expErr, err) // TODO Could/Should define error type
				}
			}
			if out != tc.expOut {
				t.Logf("Expected key: %v, got: %v", tc.expOut, out)
				t.Fail()
			}
		})
	}
}
