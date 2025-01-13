package cmd

import (
	"fmt"
	"os"

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
	Run:   runRemoveCmd,
}

func runRemoveCmd(cmd *cobra.Command, args []string) {
	internal.BlockIfNotAndroidProject(func() { os.Exit(1) })

	key := cmd.Flag("key").Value.String()
	if key == "" {
		fmt.Println("You need to pass the key through --key flag to use this command.")
		return
	}

	translations, err := internal.SingleSelectResDirectoryAndReturnTranslations()
	if err != nil || translations == nil {
		fmt.Println("Error getting translations...")
		return
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
}
