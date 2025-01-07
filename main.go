package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	_ "github.com/joho/godotenv/autoload"

	"cloud.google.com/go/translate"
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
	"google.golang.org/api/option"
)

type Translation struct {
	Path       string
	Language   string
	LocaleCode string
	RegionCode string
}

var googleApiKey = os.Getenv("GOOGLE_TRANSLATE_KEY")

func main() {
	// isAndroidProject := checkAndroidProject()
	//
	// if !isAndroidProject {
	// 	fmt.Println("This is not an Android project or you are not in the root directory of an Android project.")
	// 	os.Exit(1)
	// }
	//
	// fmt.Println("This is an Android project.")

	// currentDir, err := os.Getwd()
	// if err != nil {
	// 	fmt.Println("Error getting current directory:", err)
	// 	return
	// }

	resDirs := findResDirs("")

	if len(resDirs) == 0 {
		fmt.Println("No Android resource directories found.")
		return
	}

	fmt.Println("Android resource directories found:")
	for _, dir := range resDirs {
		fmt.Println("    ", dir)
	}
	fmt.Println("")

	workingDirectoryIndex := 0
	workingDirectory := resDirs[workingDirectoryIndex]
	fmt.Println("Working directory:", workingDirectory)

	translations := []Translation{}

	os.Chdir(workingDirectory)

	err := filepath.Walk(workingDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.Name() != "strings.xml" {
			return nil
		}

		t := getTranslationFromFileName(path)
		translations = append(translations, t)

		return nil
	})
	if err != nil {
		fmt.Println("Error walking directory:", err)
	}

	phrase := "Record Progress Stage"
	for _, translation := range translations {
		res, _ := translateText(phrase, translation.LocaleCode)
		fmt.Println(translation.Language, ":", res)
	}
}

func translateText(text, targetLanguage string) (string, error) {
	ctx := context.Background()

	client, err := translate.NewClient(ctx, option.WithAPIKey(googleApiKey))
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

func getTranslationFromFileName(path string) Translation {
	size := len(path)

	locale, region := extract(path[:size-12])
	if locale == "" {
		locale = "en"
	}

	tag := locale

	if region != "" {
		tag += "-" + region
	}

	return Translation{
		Path:       path,
		LocaleCode: locale,
		RegionCode: region,
		Language:   display.English.Languages().Name(language.MustParse(tag)),
	}
}

func extract(dirName string) (string, string) {
	fmt.Println("Extracting from:", dirName)
	re := regexp.MustCompile(`values-(\w+)(?:-r(\w+))?`)

	locale := ""
	region := ""

	matches := re.FindStringSubmatch(dirName)
	if matches != nil {
		locale = matches[1]
		region = matches[2]
	}

	return locale, region
}

func findResDirs(root string) []string {
	var resDirs []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && info.Name() == "res" {
			if isAndroidResDir(path) {
				resDirs = append(resDirs, path)
			}
		}

		return nil
	})
	if err != nil {
		fmt.Println("Error walking directory:", err)
	}

	return resDirs
}

func isAndroidResDir(path string) bool {
	if _, err := os.Stat(filepath.Join(path, "values")); !os.IsNotExist(err) {
		return true
	}

	return false
}

func checkAndroidProject() bool {
	androidRootIndicators := []string{
		"build.gradle",
		"settings.gradle",
		"settings.gradle.kts",
		"app",
	}

	isAndroidProject := false
	for _, indicator := range androidRootIndicators {
		if _, err := os.Stat(indicator); err == nil {
			isAndroidProject = true
			break
		}
	}

	return isAndroidProject
}
