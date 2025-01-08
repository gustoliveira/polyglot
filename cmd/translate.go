package cmd

import (
	"encoding/xml"
	"fmt"
	"log"

	"android-translate-tool/cmd/internal"
	"android-translate-tool/cmd/ui/singleselect"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().StringP("key", "k", "", "Key to use for translation (no spaces allowed, lowercases letters and underscores only)")
	createCmd.Flags().StringP("value", "v", "", "String to translate (english only, closed in quotes)")
	createCmd.Flags().BoolP("apply", "a", false, "Apply the translation to the project (default is false) (if false it will only print the translations)")
}

var createCmd = &cobra.Command{
	Use:   "translate",
	Short: "Translate a string",
	Run: func(cmd *cobra.Command, args []string) {
		key := cmd.Flag("key").Value.String()
		str := cmd.Flag("value").Value.String()

		fmt.Println("Key:", key)
		fmt.Println("String:", str)
		fmt.Println("")

		currentDir := "/home/gustavo/dpms/mobile-norway"
		// currentDir, err := os.Getwd()
		// if err != nil {
		// 	fmt.Println("Error getting current directory:", err)
		// 	return
		// }
		resDirs := internal.FindResourcesDirectoriesPath(currentDir)

		if len(resDirs) == 0 {
			fmt.Println("No Android resource directories found.")
			return
		}

		selectedPath := singleselect.Selection{Selected: ""}

		tprogram := tea.NewProgram(singleselect.InitialModelSingleSelect(resDirs, &selectedPath))
		if _, err := tprogram.Run(); err != nil {
			log.Printf("Name of project contains an error: %v", err)
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

		fmt.Println("Translating...")

		for _, s := range strings {
			t, err := internal.TranslateText(str, s.LocaleCode)
			if err != nil {
				fmt.Println("Error translating to", s.Language)
				continue
			}

			r, err := internal.GetResourcesFromPathXML(s.Path)
			if err != nil {
				fmt.Println(err)
				continue
			}

			r = r.AppendNewString(internal.String{
				XMLName: xml.Name{Local: "string"},
				Key:     key,
				Value:   t,
			})

			err = r.UpdateResourcesToXMLFile(s.Path)
			if err != nil {
				fmt.Println(err)
				continue
			}

			fmt.Printf("%v: %v", s.Language, t)
		}
	},
}
