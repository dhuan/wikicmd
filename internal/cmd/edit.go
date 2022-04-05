package cmd

import (
	"fmt"

	"github.com/dhuan/wikicmd/internal/utils"
	"github.com/dhuan/wikicmd/pkg/editor"
	"github.com/dhuan/wikicmd/pkg/mw"
	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit pages",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		wikiConfig, apiCredentials := beforeCommand()

		pageName := args[0]

		page, err := mw.GetPage(wikiConfig, apiCredentials, utils.FormatPageNameInput(pageName))
		if err != nil {
			fmt.Println(err)
			panic("Failed to get page.")
		}
		fmt.Println(fmt.Sprintf("Page content: %s", page.Content))

		newContent, err := editor.Edit(page.Content)
		if err != nil {
			panic("Failed to edit.")
		}
		fmt.Println(fmt.Sprintf("New content: %s", newContent))

		_, err = mw.Edit(
			wikiConfig,
			apiCredentials,
			pageName,
			newContent,
		)
		if err != nil {
			panic("Failed to edit.")
		}
	},
}
