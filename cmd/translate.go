package cmd

import (
	"encoding/xml"
	"fmt"

	"polyglot/cmd/internal"

	"github.com/spf13/cobra"
)

var (
	force     bool
	printOnly bool
)

func init() {
	rootCmd.AddCommand(translateCmd)
	translateCmd.Flags().StringP("key", "k", "", "Key to use for translation (no spaces allowed, lowercases letters and underscores only)")
	translateCmd.Flags().StringP("value", "v", "", "String to translate (english only, closed in quotes)")
	translateCmd.Flags().StringP("googleApiKey", "g", "", "Google Translate API Key (if not set it will use the GOOGLE_TRANSLATE_KEY environment variable)")
	translateCmd.Flags().BoolVar(&force, "force", false, "Force translation even if the key already exists in the file by substituting the value for the new translated one")
	translateCmd.Flags().BoolVar(&printOnly, "print-only", false, "Only print the translations, do not write to the files")
}

var translateCmd = &cobra.Command{
	Use:   "translate",
	Short: "Translate a string",
	RunE:  runTranslateCmd,
}

func runTranslateCmd(cmd *cobra.Command, args []string) error {
	err := internal.BlockIfNotAndroidProject()
	if err != nil {
		return err
	}

	key := cmd.Flag("key").Value.String()
	if !internal.IsKeyValidPrintMessage(key) {
		return fmt.Errorf("invalid key")
	}

	str := cmd.Flag("value").Value.String()
	googleApiKey := cmd.Flag("googleApiKey").Value.String()

	if googleApiKey == "" && !internal.ContainsGoogleApiKey() {
		fmt.Println("You need to pass the key through --googleApiKey flag or set the GOOGLE_TRANSLATE_KEY environment variable to use this command.")
		return fmt.Errorf("invalid googleApiKey")
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

	languagesFound := []string{}
	for _, s := range translations {
		languagesFound = append(languagesFound, s.Language)
	}

	fmt.Printf("Languages found: %v\nTranslating...\n\n", languagesFound)

	for _, t := range translations {
		r, err := internal.GetResourcesFromPathXML(t.Path)
		if err != nil {
			fmt.Println(err)
			continue
		}

		if !printOnly && (r.ContainsStringByKey(key) && !force) {
			fmt.Printf("Key <%v> already exists in %v\n", key, t.Path)
			continue
		}

		translatedText, err := internal.TranslateText(str, t.LocaleCode, &googleApiKey)
		if err != nil {
			fmt.Println("Error translating to", t.Language)
			continue
		}

		r = addStringToResources(r, t, key, translatedText)

		if !printOnly {
			err = r.UpdateResourcesToXMLFile(t.Path)
			if err != nil {
				fmt.Println(err)
				continue
			}
		}

		fmt.Printf("%v: %v\n", t.Language, translatedText)
	}

	return nil
}

func addStringToResources(r internal.Resources, t internal.Translation, key, translatedText string) internal.Resources {
	if r.ContainsStringByKey(key) && force {
		fmt.Printf("Substituting <%v> that already exists in %v\n", key, t.Path)
		return r.CreateOrSubstituteStringByKey(key, translatedText)
	}

	if r.IsSortedByKey() {
		return r.AddNewStringSorted(internal.String{
			XMLName: xml.Name{Local: "string"},
			Key:     key,
			Value:   translatedText,
		})
	}

	return r.AppendNewString(internal.String{
		XMLName: xml.Name{Local: "string"},
		Key:     key,
		Value:   translatedText,
	})
}
