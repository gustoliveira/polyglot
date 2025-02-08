package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "polyglot",
	Short: "CLI tool to manage translations in android projects using Google Translate API",
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
