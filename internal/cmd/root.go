package cmd

import (
	"github.com/spf13/cobra"
)

var (
	cfgFile     string
	userLicense string

	rootCmd = &cobra.Command{
		Use:   "wikicmd",
		Short: "Utilities for managing your Wikimedia project.",
	}
)

func Execute() error {
	rootCmd.AddCommand(editCmd)

	return rootCmd.Execute()
}
