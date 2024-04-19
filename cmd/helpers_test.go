package cmd

import (
	"testing"
)

func TestIsUrl(t *testing.T) {
	type testCase struct {
		name   string
		inURL  string
		expOut bool
	}
	tCases := []testCase{
		{
			name:   "empty string",
			inURL:  "",
			expOut: false,
		},
		{
			name:   "no scheme provided",
			inURL:  "www.example.re/",
			expOut: false,
		},
		{
			name:   "no host provided",
			inURL:  "https://",
			expOut: false,
		},
		{
			name:   "valid url provided",
			inURL:  "https://www.example.re/doc/glossary#introduction",
			expOut: true,
		},
		{
			name:   "valid url - unsupported scheme",
			inURL:  "ftp://ftp.example.re/doc/glossary.pdf",
			expOut: false,
		},
		{
			name:   "ip address",
			inURL:  "https://192.158.0.1:90",
			expOut: true,
		},
	}
	for _, tc := range tCases {
		t.Run(tc.name, func(t *testing.T) {
			result := isUrl(tc.inURL)
			if result != tc.expOut {
				t.Errorf("expected: %v, got: %v, input: %v\n", tc.expOut, result, tc.inURL)
			}
		})
	}
}
