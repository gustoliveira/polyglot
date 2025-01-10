package cmd

import (
	"fmt"
	"log"
	"os"

	"polyglot/cmd/internal"
	"polyglot/cmd/ui/singleselect"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(removeCmd)
	removeCmd.Flags().StringP("key", "k", "", "Key of the string to be removed")
}

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a key from all files of a resource directory",
	Run: func(cmd *cobra.Command, args []string) {
		key := cmd.Flag("key").Value.String()
		if key == "" {
			fmt.Println("You need to pass the key through --key flag to use this command.")
			return
		}

		currentDir, err := os.Getwd()
		if err != nil {
			fmt.Println("Error getting current directory:", err)
			return
		}

		resDirs := internal.FindResourcesDirectoriesPath(currentDir)
		if len(resDirs) == 0 {
			fmt.Println("No Android resource directories found.")
			return
		}

		selectedPath := singleselect.Selection{Selected: ""}

		tprogram := tea.NewProgram(singleselect.InitialModelSingleSelect(resDirs, &selectedPath))
		if _, err := tprogram.Run(); err != nil {
			log.Printf("Name of project contains an error: %v\n", err)
		}

		if selectedPath.Selected == "" {
			return
		}

		strings := internal.GetTranslationsFromResourceDirectory(selectedPath.Selected)

		languagesFound := []string{}
		for _, s := range strings {
			languagesFound = append(languagesFound, s.Language)
		}

		fmt.Println("Languages found:", languagesFound)

		fmt.Println("Removing...")

		for _, s := range strings {
			r, err := internal.GetResourcesFromPathXML(s.Path)
			if err != nil {
				fmt.Println(err)
				continue
			}

			if !r.ContainsStringByKey(key) {
				fmt.Printf("Key <%v> not found in %v\n", key, s.Path)
				continue
			}

			r = r.RemoveStringByKey(key)

			err = r.UpdateResourcesToXMLFile(s.Path)
			if err != nil {
				fmt.Println(err)
				continue
			}

			fmt.Printf("Removed <%v> from %v\n", key, s.Path)
		}
	},
}
