package internal

import (
	"fmt"
	"os"

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
