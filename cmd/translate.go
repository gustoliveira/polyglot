package cmd

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"

	"android-translate-tool/cmd/internal"
	"android-translate-tool/cmd/ui/singleselect"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/spf13/cobra"
)

var force bool

func init() {
	rootCmd.AddCommand(translateCmd)
	translateCmd.Flags().StringP("key", "k", "", "Key to use for translation (no spaces allowed, lowercases letters and underscores only)")
	translateCmd.Flags().StringP("value", "v", "", "String to translate (english only, closed in quotes)")
	translateCmd.Flags().StringP("googleApiKey", "g", "", "Google Translate API Key (if not set it will use the GOOGLE_TRANSLATE_KEY environment variable)")
	translateCmd.Flags().BoolVar(&force, "force", false, "Force translation even if the key already exists in the file by substituting the value for the new translated one")
}

var translateCmd = &cobra.Command{
	Use:   "translate",
	Short: "Translate a string",
	Run: func(cmd *cobra.Command, args []string) {
		key := cmd.Flag("key").Value.String()
		str := cmd.Flag("value").Value.String()
		googleApiKey := cmd.Flag("googleApiKey").Value.String()

		if googleApiKey == "" && !internal.ContainsGoogleApiKey() {
			fmt.Println("You need to pass the key through --googleApiKey flag or set the GOOGLE_TRANSLATE_KEY environment variable to use this command.")
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

		fmt.Printf("Languages found: %v\n", languagesFound)

		fmt.Println("Translating...\n")

		for _, s := range strings {
			r, err := internal.GetResourcesFromPathXML(s.Path)
			if err != nil {
				fmt.Println(err)
				continue
			}

			if r.ContainsStringByKey(key) && !force {
				fmt.Printf("Key <%v> already exists in %v\n", key, s.Path)
				continue
			}

			t, err := internal.TranslateText(str, s.LocaleCode, &googleApiKey)
			if err != nil {
				fmt.Println("Error translating to", s.Language)
				continue
			}

			if r.ContainsStringByKey(key) && force {
				fmt.Printf("Substituting <%v> that already exists in %v\n", key, s.Path)
				r = r.CreateOrSubstituteStringByKey(key, t)
			} else {
				r = r.AppendNewString(internal.String{
					XMLName: xml.Name{Local: "string"},
					Key:     key,
					Value:   t,
				})
			}

			err = r.UpdateResourcesToXMLFile(s.Path)
			if err != nil {
				fmt.Println(err)
				continue
			}

			fmt.Printf("%v: %v\n", s.Language, t)
		}
	},
}
