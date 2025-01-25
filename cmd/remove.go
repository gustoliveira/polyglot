package cmd

import (
	"fmt"

	"polyglot/cmd/internal"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(removeCmd)
	removeCmd.Flags().StringP("key", "k", "", "Key of the string to be removed")
}

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a key from all files of a resource directory",
	RunE:  runRemoveCmd,
}

func runRemoveCmd(cmd *cobra.Command, args []string) error {
	err := internal.BlockIfNotAndroidProject()
	if err != nil {
		return err
	}

	key := cmd.Flag("key").Value.String()
	if !internal.IsKeyValidPrintMessage(key) {
		return fmt.Errorf("invalid key")
	}

	translations, err := internal.SingleSelectResDirectoryAndReturnTranslations()
	if err != nil || translations == nil {
		if err != nil {
			return err
		}
		if translations != nil {
			return fmt.Errorf("no translations found")
		}
	}

	fmt.Println("Removing...")

	for _, t := range translations {
		r, err := internal.GetResourcesFromPathXML(t.Path)
		if err != nil {
			fmt.Println(err)
			continue
		}

		if !r.ContainsStringByKey(key) {
			fmt.Printf("Key <%v> not found in %v\n", key, t.Path)
			continue
		}

		r = r.RemoveStringByKey(key)

		err = r.UpdateResourcesToXMLFile(t.Path)
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Printf("Removed <%v> from %v\n", key, t.Path)
	}

	return nil
}
