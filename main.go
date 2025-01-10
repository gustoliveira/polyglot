package main

import (
	"fmt"
	"os"

	"polyglot/cmd"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	inAndroidProject := CheckCurrentDirectoryIsAndroidProject()
	if !inAndroidProject {
		fmt.Println("This is not an Android project or you are not in the root directory of an Android project.")
		return
	}

	cmd.Execute()
}

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
