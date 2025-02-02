package cmd

import (
	"fmt"

	"polyglot/cmd/internal"

	"github.com/spf13/cobra"
)

var allModulesC bool

func init() {
	// TODO: Add flags to select which type of normalization to check
	// If none selected check all
	rootCmd.AddCommand(checkCmd)
	checkCmd.Flags().BoolVar(&allModulesC, "all", false, "Check all translations files of all project modules")
}

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check if the translations is properly normalized",
	RunE:  runCheckCmd,
}

func runCheckCmd(cmd *cobra.Command, args []string) error {
	err := internal.BlockIfNotAndroidProject()
	if err != nil {
		return err
	}

	translations, err := internal.GetTranslations(allModulesC)
	if err != nil || translations == nil {
		if err != nil {
			return err
		}
		if translations != nil {
			return fmt.Errorf("no translations found")
		}
	}

	allResources := internal.ListResources{}

	keys := make(map[string]struct{})

	// CHECK: Traslations files are sorted by key
	fmt.Printf("Checking if translation files are sorted by key...\n")
	for _, t := range translations {
		r, err := internal.GetResourcesFromPathXML(t.Path)
		if err != nil {
			fmt.Println(err)
			continue
		}

		allResources = append(allResources, r)

		if !r.IsSortedByKey() {
			fmt.Printf("\tFAIL: File \"%v\" is not sorted by key\n", t.Path)
		} else {
			fmt.Printf("\tPASS: File \"%v\" is sorted by key\n", t.Path)
		}

		// Add all keys of all resources to check for possible unused keys
		for _, key := range r.Strings {
			// Set is used as a map with no values
			keys[key.Key] = struct{}{}
		}
	}

	// CHECK: Find possible unused keys
	fmt.Printf("\nSearch for *possible* unused keys in all resources...\n")
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
	fmt.Printf("Found %v possible unused keys\n\n", count)

	// CHECK: Missing translations between files
	fmt.Printf("Checking for *possible* missing translations between files...\n")
	missingTranslationRelatory := allResources.CheckMissingTranslations().CheckMissingTranslationsRelatory()
	fmt.Println(missingTranslationRelatory)

	return nil
}
