package internal

import (
	"reflect"
	"testing"
)

func TestIsStringKeyValid(t *testing.T) {
	testCases := []struct {
		name     string
		key      string
		response bool
	}{
		{name: "Valid string key 1", key: "palmeiras", response: true},
		{name: "Valid string key 2", key: "sociedade_esportiva_palmeiras", response: true},
		{name: "Valid string key 3", key: "p", response: true},
		{name: "Invalid string key 1", key: "", response: false},
		{name: "Invalid string key 2", key: "palmeiras ", response: false},
		{name: "Invalid string key 3", key: "1914", response: false},
		{name: "Invalid string key 4", key: "palmeiras1914", response: false},
		{name: "Invalid string key 4", key: "sociedade esportiva palmeiras", response: false},
		{name: "Invalid string key 5", key: "sociedade-esportiva-palmeiras", response: false},
		{name: "Invalid string key 6", key: "palmeiras!", response: false},
		{name: "Invalid string key 7", key: "Palmeiras", response: false},
		{name: "Invalid string key 8", key: "_", response: false},
		{name: "Invalid string key 9", key: "_palmeiras", response: false},
		{name: "Invalid string key 10", key: "palmeiras_", response: false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := IsStringKeyValid(tc.key)

			if !reflect.DeepEqual(result, tc.response) {
				t.Errorf("IsStringKeyValid() = %v, want %v", result, tc.response)
			}
		})
	}
}

func TestIsKeyValidPrintMessage(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedValid bool
		expectedMsg   string
	}{
		{
			name:          "Empty Key",
			input:         "",
			expectedValid: false,
			expectedMsg:   "You need to pass the key through --key flag to use this command.",
		},
		{
			name:          "Invalid Key - starts with underscore",
			input:         "_dev",
			expectedValid: false,
			expectedMsg:   "Invalid key. Only lowercases letters and underscores are allowed.",
		},
		{
			name:          "Invalid Key - ends with underscore",
			input:         "dev_",
			expectedValid: false,
			expectedMsg:   "Invalid key. Only lowercases letters and underscores are allowed.",
		},
		{
			name:          "Invalid Key - no lowercase letters",
			input:         "_",
			expectedValid: false,
			expectedMsg:   "Invalid key. Only lowercases letters and underscores are allowed.",
		},
		{
			name:          "Valid Key - single character",
			input:         "d",
			expectedValid: true,
			expectedMsg:   "",
		},
		{
			name:          "Valid Key - multiple lowercase letters",
			input:         "palmeiras",
			expectedValid: true,
			expectedMsg:   "",
		},
		{
			name:          "Valid Key - with underscore",
			input:         "palmeiras_teste",
			expectedValid: true,
			expectedMsg:   "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := IsKeyValidPrintMessage(tc.input)
			if result != tc.expectedValid {
				t.Errorf("IsKeyValidPrintMessage(%q) = %v; want %v", tc.input, result, tc.expectedValid)
			}

			output := CaptureStdout(func() {
				IsKeyValidPrintMessage(tc.input)
			})
			if tc.expectedMsg != "" && output != tc.expectedMsg+"\n" {
				t.Errorf("Expected message %q, got %q", tc.expectedMsg, output)
			}

			if tc.expectedMsg == "" && output != "" {
				t.Errorf("Did not expect any message, but got %q", output)
			}
		})
	}
}
