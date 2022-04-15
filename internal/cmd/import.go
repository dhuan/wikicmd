package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/dhuan/wikicmd/internal/utils"
	"github.com/dhuan/wikicmd/pkg/mw"
	"github.com/spf13/cobra"
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import pages and images",
	Run: func(cmd *cobra.Command, filePaths []string) {
		wikiConfig, apiCredentials, hook := beforeCommand()

		fileValidationErrors := utils.ValidateFiles(filePaths, allowedExtensionsToBeImported)
		if len(fileValidationErrors) > 0 {
			handleFileValidationErrors(fileValidationErrors)

			os.Exit(1)
		}

		uploadedCount := runImport(wikiConfig, apiCredentials, filePaths, hook)

		fmt.Println(fmt.Sprintf("%d item(s) have been imported.\nDone!", uploadedCount))
	},
}

func runImport(
	wikiConfig *mw.Config,
	apiCredentials *mw.ApiCredentials,
	filePaths []string,
	hook *mw.HookOptions,
) int {
	uploadedCount := 0
	for _, filePath := range filePaths {
		fmt.Println(fmt.Sprintf("Importing %s", filePath))

		file, err := os.Open(filePath)
		if err != nil {
			panic(err)
		}

		fileContent, err := os.ReadFile(filePath)
		if err != nil {
			panic(err)
		}

		if fileIsPage(filePath) {
			if err = importPage(wikiConfig, apiCredentials, filePath, fileContent, hook); err != nil {
				panic(err)
			}

			uploadedCount = uploadedCount + 1

			continue
		}

		if fileIsImage(filePath) {
			uploadWarnings, uploaded, err := importImage(wikiConfig, apiCredentials, filePath, file)
			if err != nil {
				panic(err)
			}

			handleUploadWarnings(uploadWarnings)

			if uploaded {
				uploadedCount = uploadedCount + 1
			}

			continue
		}

		panic("Something went wrong. Could not resolve type of file to be imported.")
	}

	return uploadedCount
}

func fileIsPage(filePath string) bool {
	return utils.ExtensionMatches(pageExtensions, filePath)
}

func fileIsImage(filePath string) bool {
	return utils.ExtensionMatches(imageExtensions, filePath)
}

func importImage(
	wikiConfig *mw.Config,
	apiCredentials *mw.ApiCredentials,
	filePath string,
	file io.Reader,
) ([]mw.UploadWarning, bool, error) {
	fileName := filepath.Base(filePath)

	err, warnings, uploaded := mw.Upload(wikiConfig, apiCredentials, fileName, file)
	if err != nil {
		return []mw.UploadWarning{}, uploaded, err
	}

	return warnings, uploaded, nil
}

func importPage(
	wikiConfig *mw.Config,
	apiCredentials *mw.ApiCredentials,
	filePath string,
	fileContent []byte,
	hook *mw.HookOptions,
) error {
	pageName := utils.FilePathToPageName(filePath)
	_, err := mw.Edit(
		wikiConfig,
		apiCredentials,
		pageName,
		string(fileContent),
		hook,
	)
	if err != nil {
		return err
	}

	return nil
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

func handleUploadWarnings(warnings []mw.UploadWarning) {
	if len(warnings) == 0 {
		return
	}

	for _, warning := range warnings {
		message, ok := MAP_UPLOAD_WARNING_MESSAGE[warning]

		if ok {
			fmt.Println(fmt.Sprintf("WARNING: %s", message))
		}
	}
}
