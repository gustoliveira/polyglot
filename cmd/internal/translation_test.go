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

func TestFormatTranslation(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"NoSingleQuotes", `Hello world`, `Hello world`},
		{"RealUseCase", `Aucun élément trouvé avec les filtres actuels. Essayez de modifier les filtres pour afficher plus d'éléments.`, `Aucun élément trouvé avec les filtres actuels. Essayez de modifier les filtres pour afficher plus d\'éléments.`},
		{"SingleQuoteAtBeginning", `'Hello world`, `\'Hello world`},
		{"SingleQuoteAtEnd", `Hello world'`, `Hello world\'`},
		{"SingleQuoteInMiddle", `Hello 'world'`, `Hello \'world\'`},
		{"MultipleSingleQuotes", `'It's a beautiful day'`, `\'It\'s a beautiful day\'`},
		{"OnlySingleQuotes", `'''`, `\'\'\'`},
		{"EmptyString", ``, ``},
		{"MixedQuotes", `She said 'Hello "world"'`, `She said \'Hello "world"\'`},
		{"ConsecutiveSingleQuotes", `It''s working`, `It\'\'s working`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatTranslation(tt.input)
			if result != tt.expected {
				t.Errorf("FormatTranslation(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
