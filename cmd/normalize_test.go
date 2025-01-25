package cmd

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeCmd_non_android_project_directory(t *testing.T) {
	root := &cobra.Command{Use: "normalize", RunE: normalizeCmd.RunE}
	err := root.Execute()

	assert.Error(t, err)
	assert.Equal(t, "current directory is not an android project", err.Error())
}
