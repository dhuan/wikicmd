package cmd

import (
	"github.com/spf13/cobra"
)

var (
	//Flags
	flagVerbose     bool
	flagConfigNew   bool
	flagExportType  string
	flagEditMessage string

	rootCmd = &cobra.Command{
		Use:   "wikicmd",
		Short: "Utilities for managing your MediaWiki project.",
	}
)

func Execute() error {
	rootCmd.PersistentFlags().BoolVarP(&flagVerbose, "verbose", "v", false, "Verbose output.")
	configCmd.Flags().BoolVarP(&flagConfigNew, "new", "n", false, "Create new configuration file even if one already exists.")
	exportCmd.Flags().StringVarP(&flagExportType, "type", "t", "all", "Which type of item to export. (all/page/image)")
	editCmd.Flags().StringVarP(&flagEditMessage, "message", "m", "", "Summary for your change.")

	rootCmd.AddCommand(editCmd)
	rootCmd.AddCommand(exportCmd)
	rootCmd.AddCommand(importCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(switchCmd)
	rootCmd.AddCommand(versionCmd)

	return rootCmd.Execute()
}
