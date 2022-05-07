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
		wikiConfig, apiCredentials, hookOptions := beforeCommand(true)
		pageName := args[0]

		page, err := mw.GetPage(
			wikiConfig,
			apiCredentials,
			utils.FormatPageNameInput(pageName),
			hookOptions,
		)
		if err != nil {
			fmt.Println(err)
			panic("Failed to get page.")
		}

		newContent, changed, err := editor.Edit(page.Content)
		if err != nil {
			panic("Failed to edit.")
		}

		if changed {
			_, err = mw.Edit(
				wikiConfig,
				apiCredentials,
				pageName,
				newContent,
				hookOptions,
			)
			if err != nil {
				panic("Failed to edit.")
			}

			if page.Exists {
				fmt.Println(fmt.Sprintf("%s edited successfully.", page.Name))
			} else {
				fmt.Println(fmt.Sprintf("%s created successfully.", page.Name))
			}

			return
		}

		fmt.Println("No changes were made.")
	},
}
