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

	keys := make(map[string]struct{})

	for _, t := range translations {
		r, err := internal.GetResourcesFromPathXML(t.Path)
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Printf("Checking if translation files are sorted by key in %v...\n", t.Path)
		if !r.IsSortedByKey() {
			fmt.Printf("\tFAIL: File %v is not sorted by key\n", t.Path)
		}

		// Add all keys of all resources to check for possible unused keys
		for _, key := range r.Strings {
			// Set is used as a map with no values
			keys[key.Key] = struct{}{}
		}
	}

	fmt.Printf("Search for *possible* unused keys in all resources...\n")
	count := 0
	for key := range keys {
		isBeingUsed, err := internal.IsKeyBeingUsed(key)
		if err != nil {
			fmt.Println(err)
			continue
		}

		if !isBeingUsed {
			fmt.Printf("\tString <%v> appears to be unused\n", fmt.Sprintf("R.string.%v", key))
			count++
		}
	}

	fmt.Printf("Found %v possible unused keys\n\n\n", count)
}
