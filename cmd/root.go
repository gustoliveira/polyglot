package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "polyglot",
	Short: "A CLI tool to translate new translations in an android project and insert automatically using Google Translate API",
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
