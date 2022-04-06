package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/dhuan/wikicmd/internal/utils"
	"github.com/dhuan/wikicmd/pkg/mw"
	"github.com/spf13/cobra"
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import pages",
	Run: func(cmd *cobra.Command, filePaths []string) {
		wikiConfig, apiCredentials := beforeCommand()

		fileValidationErrors := utils.ValidateFiles(filePaths, []string{"wikitext"})
		if len(fileValidationErrors) > 0 {
			handleFileValidationErrors(fileValidationErrors)

			os.Exit(1)
		}

		runImport(wikiConfig, apiCredentials, filePaths)

		fmt.Println(fmt.Sprintf("%d page(s) have been imported.\nDone!", len(filePaths)))
	},
}

func runImport(wikiConfig *mw.Config, apiCredentials *mw.ApiCredentials, filePaths []string) {
	for _, filePath := range filePaths {
		fmt.Println(fmt.Sprintf("Importing %s", filePath))

		fileContent, err := os.ReadFile(filePath)
		if err != nil {
			panic("Failed to read file.")
		}

		pageName := utils.FilePathToPageName(filePath)
		_, err = mw.Edit(
			wikiConfig,
			apiCredentials,
			pageName,
			string(fileContent),
		)
		if err != nil {
			panic("Failed to edit.")
		}
	}
}

func handleFileValidationErrors(fileValidationErrors map[string]error) {
	fmt.Println("The following errors occurred:")

	for filePath, fileValidationError := range fileValidationErrors {
		if errors.Is(fileValidationError, utils.ErrFileDoesNotExist) {
			fmt.Println(fmt.Sprintf("%s: Does not exist.", filePath))

			continue
		}

		if errors.Is(fileValidationError, utils.ErrExtensionNotAccepted) {
			fmt.Println(fmt.Sprintf("%s: Extension not allowed.", filePath))

			continue
		}

		panic("This error is unknown.")
	}
}
