package cmd

import (
	"fmt"
	"os"

	"polyglot/cmd/internal"

	"github.com/spf13/cobra"
)

func init() {
	// TODO: Add flags to select which type of normalization to apply
	// If none selected apply all
	rootCmd.AddCommand(normalizeCmd)
}

var normalizeCmd = &cobra.Command{
	Use:   "normalize",
	Short: "Normalize translations files",
	Run:   runNormalizeCmd,
}

func runNormalizeCmd(cmd *cobra.Command, args []string) {
	internal.BlockIfNotAndroidProject(func() { os.Exit(1) })

	translations, err := internal.SingleSelectResDirectoryAndReturnTranslations()
	if err != nil || translations == nil {
		fmt.Println("Error getting translations...")
		return
	}

	for _, t := range translations {
		r, err := internal.GetResourcesFromPathXML(t.Path)
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Println("Sorting by key unsorted translation files...")
		if !r.IsSortedByKey() {
			fmt.Printf("File %v is not sorted by key. Sorting...\n", t.Path)

			r.SortByKey()

			err = r.UpdateResourcesToXMLFile(t.Path)
			if err != nil {
				fmt.Println(err)
				continue
			}
		}
	}
}
