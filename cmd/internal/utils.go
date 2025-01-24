package internal

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"

	"polyglot/cmd/ui/singleselect"

	tea "github.com/charmbracelet/bubbletea"
)

func SingleSelectResDirectoryAndReturnTranslations() ([]Translation, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return nil, err
	}

	resDirs := FindResourcesDirectoriesPath(currentDir)
	if len(resDirs) == 0 {
		fmt.Println("No Android resource directories found.")
		return nil, err
	}

	selectedPath := singleselect.Selection{Selected: ""}

	tprogram := tea.NewProgram(singleselect.InitialModelSingleSelect(resDirs, &selectedPath))
	if _, err := tprogram.Run(); err != nil {
		fmt.Printf("Name of project contains an error: %v\n", err)
		return nil, err
	}

	if selectedPath.Selected == "" {
		return nil, nil
	}

	translations := GetTranslationsFromResourceDirectory(selectedPath.Selected)

	return translations, nil
}

func IsKeyBeingUsed(key string) (bool, error) {
	_, err := exec.LookPath("rg")
	if err == nil {
		return IsKeyBeingUsedRipGrep(key)
	}

	return IsKeyBeingUsedGrep(key)
}

func IsKeyBeingUsedRipGrep(key string) (bool, error) {
	k := fmt.Sprintf("R.string.%v", key)

	cmd := exec.Command("rg", k, "--glob=*.kt")

	var output bytes.Buffer
	cmd.Stderr = &output

	err := cmd.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if exitErr.ExitCode() == 1 {
				return false, nil
			}
		}

		return false, err
	}

	return true, nil
}

func IsKeyBeingUsedGrep(key string) (bool, error) {
	k := fmt.Sprintf("R.string.%v", key)

	cmd := exec.Command("grep", "-r", k, "--include=*.kt", ".")

	var output bytes.Buffer
	cmd.Stderr = &output

	err := cmd.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if exitErr.ExitCode() == 1 {
				return false, nil
			}
		}

		return false, err
	}

	return true, nil
}

func IsStringKeyValid(k string) bool {
	isValid := regexp.MustCompile(`^[a-z](?:[a-z_]*[a-z])*$`).MatchString
	return len(k) != 0 && isValid(k)
}

func IsKeyValidPrintMessage(key string) bool {
	if key == "" {
		fmt.Println("You need to pass the key through --key flag to use this command.")
		return false
	}
	if !IsStringKeyValid(key) {
		fmt.Println("Invalid key. Only lowercases letters and underscores are allowed.")
		return false
	}

	return true
}
