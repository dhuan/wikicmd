package cmd

import (
	"fmt"

	"github.com/dhuan/wikicmd/internal/wikicmd"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print current version.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(wikicmd.CurrentVersion())
	},
}
