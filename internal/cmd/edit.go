package cmd

import (
	"fmt"

	"github.com/dhuan/wikicmd/internal/input"
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
			summary, err := promptSummary()
			if err != nil {
				panic(err)
			}

			_, err = mw.Edit(
				wikiConfig,
				apiCredentials,
				pageName,
				newContent,
				summary,
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

func promptSummary() (string, error) {
	inputSummary, err := input.ResponsePrompt[bool](
		"Would you like to set a summary for this change (yes/no): ",
		map[string]bool{
			"yes": true,
			"y":   true,
			"no":  false,
		},
		false,
		false,
	)
	if err != nil {
		return "", err
	}

	if !inputSummary {
		return "", nil
	}

	summaryContent, _, err := editor.Edit("")
	if err != nil {
		return "", err
	}

	return summaryContent, nil
}
