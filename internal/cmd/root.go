package cmd

import (
	"github.com/spf13/cobra"
)

var (
	//Flags
	flagVerbose    bool
	flagConfigNew  bool
	flagExportType string

	rootCmd = &cobra.Command{
		Use:   "wikicmd",
		Short: "Utilities for managing your Wikimedia project.",
	}
)

func Execute() error {
	rootCmd.PersistentFlags().BoolVarP(&flagVerbose, "verbose", "v", false, "Verbose output.")
	configCmd.Flags().BoolVarP(&flagConfigNew, "new", "n", false, "Create new configuration file even if one already exists.")
	exportCmd.Flags().StringVarP(&flagExportType, "type", "t", "all", "Which type of item to export. (all/page/image)")

	rootCmd.AddCommand(editCmd)
	rootCmd.AddCommand(exportCmd)
	rootCmd.AddCommand(importCmd)
	rootCmd.AddCommand(configCmd)

	return rootCmd.Execute()
}
