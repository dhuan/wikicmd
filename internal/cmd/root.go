package cmd

import (
	"github.com/spf13/cobra"
)

var (
	cfgFile     string
	userLicense string

	//Flags
	FlagConfigNew bool

	rootCmd = &cobra.Command{
		Use:   "wikicmd",
		Short: "Utilities for managing your Wikimedia project.",
	}
)

func Execute() error {
	configCmd.Flags().BoolVarP(&FlagConfigNew, "new", "n", false, "Create new configuration file even if one already exists.")

	rootCmd.AddCommand(editCmd)
	rootCmd.AddCommand(exportCmd)
	rootCmd.AddCommand(importCmd)
	rootCmd.AddCommand(configCmd)

	return rootCmd.Execute()
}
