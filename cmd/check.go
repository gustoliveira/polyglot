package cmd

import (
	"fmt"
	"os"

	"polyglot/cmd/internal"

	"github.com/spf13/cobra"
)

func init() {
	// TODO: Add flags to select which type of normalization to check
	// If none selected check all
	rootCmd.AddCommand(checkCmd)
}

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check if the translations is properly normalized",
	Run:   runCheckCmd,
}

func runCheckCmd(cmd *cobra.Command, args []string) {
	internal.BlockIfNotAndroidProject(func() { os.Exit(1) })

	translations, err := internal.SingleSelectResDirectoryAndReturnTranslations()
	if err != nil || translations == nil {
		fmt.Println("Error getting translations...")
		return
	}

	fmt.Println("Checking if translation files are sorted by key...")

	for _, t := range translations {
		r, err := internal.GetResourcesFromPathXML(t.Path)
		if err != nil {
			fmt.Println(err)
			continue
		}

		if !r.IsSortedByKey() {
			fmt.Printf("FAIL: File %v is not sorted by key\n", t.Path)
		}
	}
}
