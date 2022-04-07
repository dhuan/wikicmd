package cmd

import (
	"github.com/dhuan/wikicmd/internal/config"
	"github.com/dhuan/wikicmd/pkg/editor"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Edit your configuration file.",
	Run: func(cmd *cobra.Command, filePaths []string) {
		configFilePath, _, err := config.GetConfigFilePath()
		if err != nil {
			panic(err)
		}

		err = editor.EditFile(configFilePath)
		if err != nil {
			panic(err)
		}
	},
}
