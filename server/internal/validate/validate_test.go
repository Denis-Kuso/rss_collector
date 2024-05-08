package validate

import (
	"testing"
)

func TestIsValidUUID(t *testing.T) {
	type testCase struct {
		name   string
		input  string
		expOut bool
	}
	testCases := []testCase{
		{
			name:   "valid UUID",
			input:  "c5c9212c-57a3-4d68-b42e-addd951502c0",
			expOut: false,
		},
		{
			name:   "off by 1 UUID",
			input:  "c5c9212c-57a3-4d68-b42e-addd951502c",
			expOut: false,
		},
		{
			name:   "blank",
			input:  "",
			expOut: false,
		},
		{
			name:   "unrelated string",
			input:  "fkolaspkrwpe0-lj",
			expOut: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := ValidLimit(tc.input)
			if result != tc.expOut {
				t.Errorf("input: \"%s\", expected: %t, got %t\n", tc.input, tc.expOut, result)
			}
		})
	}
}

func TestValidLimit(t *testing.T) {
	type testCase struct {
		name   string
		input  string
		expOut bool
	}
	testCases := []testCase{
		{
			name:   "happy case",
			input:  "1",
			expOut: true,
		},
		{
			name:   "invalid input - negative num",
			input:  "-2",
			expOut: false,
		},
		{
			name:   "invalid input - NaN",
			input:  "plane",
			expOut: false,
		},
		{
			name:   "max limit exceeded",
			input:  "1000",
			expOut: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := ValidLimit(tc.input)
			if result != tc.expOut {
				t.Errorf("input: \"%s\", expected: %t, got %t\n", tc.input, tc.expOut, result)
			}
		})
	}
}

func TestValidateUsername(t *testing.T) {
	const MAX_LENGTH = 35
	type testCase struct {
		name   string
		input  string
		expOut bool
	}
	tCases := []testCase{
		{
			name:   "empty string",
			input:  "",
			expOut: false,
		},
		{
			name:   "invalid length",
			input:  "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
			expOut: false,
		},
		{
			name:   "max allowed length",
			input:  "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", // MAX_LENGTH,
			expOut: true,
		},
		{
			name:   "one character",
			input:  "£",
			expOut: true,
		},
	}
	for _, tc := range tCases {
		t.Run(tc.name, func(t *testing.T) {
			result := ValidateUsername(tc.input)
			if result != tc.expOut {
				t.Errorf("input: \"%s\", expected: %t, got %t\n", tc.input, tc.expOut, result)
			}
		})
	}
}

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
			result := IsUrl(tc.inURL)
			if result != tc.expOut {
				t.Errorf("expected: %v, got: %v, input: %v\n", tc.expOut, result, tc.inURL)
			}
		})
	}
}
