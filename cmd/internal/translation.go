package internal

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/translate"
	"golang.org/x/text/language"
	"google.golang.org/api/option"
)

var GOOGLE_API_KEY = os.Getenv("GOOGLE_TRANSLATE_KEY")

type Translation struct {
	Path       string
	Language   string
	LocaleCode string
	RegionCode string
}

func ContainsGoogleApiKey() bool {
	return GOOGLE_API_KEY != ""
}

func TranslateText(text string, targetLanguage string, googleApiKey *string) (string, error) {
	key := GOOGLE_API_KEY
	if googleApiKey != nil && *googleApiKey != "" {
		key = *googleApiKey
	}

	ctx := context.Background()

	client, err := translate.NewClient(ctx, option.WithAPIKey(key))
	if err != nil {
		return "", fmt.Errorf("failed to create client: %v", err)
	}
	defer client.Close()

	lang, err := language.Parse(targetLanguage)
	if err != nil {
		return "", fmt.Errorf("failed to parse target language: %v", err)
	}

	resp, err := client.Translate(ctx, []string{text}, lang, nil)
	if err != nil {
		return "", fmt.Errorf("failed to translate text: %v", err)
	}

	if len(resp) == 0 {
		return "", fmt.Errorf("translation response is empty")
	}

	return resp[0].Text, nil
}
