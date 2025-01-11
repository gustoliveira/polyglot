package internal

import (
	"os"
	"testing"
)

func resetGoogleAPIKey() {
	GOOGLE_API_KEY = os.Getenv("GOOGLE_TRANSLATE_KEY")
}

func TestContainsGoogleApiKey(t *testing.T) {
	tests := []struct {
		name   string
		envVal string
		want   bool
	}{
		{"EmptyVariable", "", false},
		{"HasKey", "YOUR_API_KEY", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("GOOGLE_TRANSLATE_KEY", tt.envVal)
			defer os.Unsetenv("GOOGLE_TRANSLATE_KEY")
			resetGoogleAPIKey()

			got := ContainsGoogleApiKey()
			if got != tt.want {
				t.Errorf("ContainsGoogleApiKey() = %v, want %v", got, tt.want)
			}
		})
	}

	t.Run("Test without set any variable", func(t *testing.T) {
		os.Unsetenv("GOOGLE_TRANSLATE_KEY")
		resetGoogleAPIKey()

		got := ContainsGoogleApiKey()
		if got != false {
			t.Errorf("ContainsGoogleApiKey() = %v, want %v", got, false)
		}
	})
}
