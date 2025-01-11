package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
)

func CheckCurrentDirectoryIsAndroidProject() bool {
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

func BlockIfNotAndroidProject(onBlockCallback func()) {
	if !CheckCurrentDirectoryIsAndroidProject() {
		fmt.Println("Current directory is not an Android project.")
		onBlockCallback()
	}
}

func FindResourcesDirectoriesPath(root string) []string {
	var resDirs []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && info.Name() == "res" {
			if isAndroidResourceDirectory(path) {
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

func isAndroidResourceDirectory(path string) bool {
	if _, err := os.Stat(filepath.Join(path, "values")); !os.IsNotExist(err) {
		return true
	}

	return false
}

func GetTranslationsFromResourceDirectory(path string) []Translation {
	translations := []Translation{}

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.Name() != "strings.xml" {
			return nil
		}

		t, err := GetTranslationFromFileName(path)
		if err != nil {
			return err
		}

		translations = append(translations, t)

		return nil
	})
	if err != nil {
		fmt.Println("Error walking directory:", err)
	}

	return translations
}

func GetTranslationFromFileName(path string) (Translation, error) {
	size := len(path)

	locale, region := extract(path[:size-12])
	if locale == "" {
		locale = "en"
	}

	languageTag := locale

	if region != "" {
		languageTag += "-" + region
	}

	tag, err := language.Parse(languageTag)
	if err != nil {
		fmt.Println("Error parsing language tag:", err)
		return Translation{}, err
	}

	return Translation{
		Path:       path,
		LocaleCode: locale,
		RegionCode: region,
		Language:   display.English.Languages().Name(tag),
	}, nil
}

func extract(dirName string) (string, string) {
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
