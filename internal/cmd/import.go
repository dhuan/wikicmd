package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/dhuan/wikicmd/internal/config"
	"github.com/dhuan/wikicmd/internal/utils"
	"github.com/dhuan/wikicmd/pkg/mw"
	"github.com/spf13/cobra"
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import pages and images",
	Run: func(cmd *cobra.Command, filePaths []string) {
		wikiConfig, apiCredentials, _, hook := beforeCommand(true)
		userConfig, _, err := config.Get()
		if err != nil {
			panic(err)
		}

		allAllowedExtensions := append(
			config.ImportExtensionsPage(),
			config.ImportExtensionsMedia(userConfig)...,
		)
		fileValidationErrors := utils.ValidateFiles(filePaths, allAllowedExtensions)
		if len(fileValidationErrors) > 0 {
			handleFileValidationErrors(fileValidationErrors)

			os.Exit(1)
		}

		uploadedCount := runImport(userConfig, wikiConfig, apiCredentials, filePaths, hook)

		fmt.Printf("%d item(s) have been imported.\nDone!\n", uploadedCount)
	},
}

func runImport(
	userConfig *config.WikiConfig,
	wikiConfig *mw.Config,
	apiCredentials *mw.ApiCredentials,
	filePaths []string,
	hook *mw.HookOptions,
) int {
	uploadedCount := 0
	for _, filePath := range filePaths {
		fmt.Printf("Importing %s\n", filePath)

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

		if fileIsImage(userConfig, filePath) {
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
	return utils.ExtensionMatches(config.ImportExtensionsPage(), filePath)
}

func fileIsImage(userConfig *config.WikiConfig, filePath string) bool {
	return utils.ExtensionMatches(config.ImportExtensionsMedia(userConfig), filePath)
}

func importImage(
	wikiConfig *mw.Config,
	apiCredentials *mw.ApiCredentials,
	filePath string,
	file io.Reader,
) ([]mw.UploadWarning, bool, error) {
	fileName := filepath.Base(filePath)

	warnings, uploaded, err := mw.Upload(wikiConfig, apiCredentials, fileName, file)
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
	pageName := utils.FilePathToPageName(config.ImportExtensionsPage(), filePath)
	_, err := mw.Edit(
		wikiConfig,
		apiCredentials,
		pageName,
		string(fileContent),
		flagMessage,
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
			fmt.Printf("%s: Does not exist.\n", filePath)

			continue
		}

		if errors.Is(fileValidationError, utils.ErrExtensionNotAccepted) {
			fmt.Printf("%s: Extension not allowed.\n", filePath)

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
			fmt.Printf("WARNING: %s\n", message)
		}
	}
}
