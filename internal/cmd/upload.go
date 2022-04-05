package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dhuan/wikicmd/internal/utils"
	"github.com/dhuan/wikicmd/pkg/mw"
	"github.com/spf13/cobra"
)

var MAP_WARNING_MESSAGE = map[mw.UploadWarning]string{
	mw.UPLOAD_WARNING_SAME_FILE_NO_CHANGE: "File was not uploaded because the existing image is exactly the same.",
}

var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload images",
	Run: func(cmd *cobra.Command, filePaths []string) {
		wikiConfig, apiCredentials := beforeCommand()

		failedImages := validateImages(filePaths)
		if len(failedImages) > 0 {
			failedImagesStr := strings.Join(failedImages, "\n")
			fmt.Println((fmt.Sprintf("The following files cannot be uploaded:\n%s", failedImagesStr)))

			os.Exit(1)
		}

		uploadedCount := 0
		for _, filePath := range filePaths {
			fmt.Println(fmt.Sprintf("Uploading %s", filePath))

			fileContent, err := os.Open(filePath)
			if err != nil {
				panic(fmt.Sprintf("Failed to upload: %s", filePath))
			}

			fileName := filepath.Base(filePath)

			err, warnings, uploaded := mw.Upload(wikiConfig, apiCredentials, fileName, fileContent)
			if uploaded {
				uploadedCount = uploadedCount + 1
			}
			if err != nil {
				panic(err)
			}

			if len(warnings) > 0 {
				handleUploadWarnings(warnings)
			} else {
				fmt.Println(fmt.Sprintf("File uploaded successfully: %s.", fileName))
			}
		}

		fmt.Println(fmt.Sprintf("%d file(s) have been uploaded.\nDone!", uploadedCount))
	},
}

func handleUploadWarnings(warnings []mw.UploadWarning) {
	for _, warning := range warnings {
		message, ok := MAP_WARNING_MESSAGE[warning]

		if ok {
			fmt.Println(fmt.Sprintf("WARNING: %s", message))
		}
	}
}

func validateImages(filePaths []string) []string {
	failedImages := make([]string, 0)

	for _, filePath := range filePaths {
		if extensionIsValid := utils.ExtensionMatches(allowedExtensionsToBeUploaded, filePath); !extensionIsValid {
			failedImages = append(failedImages, filePath)

			continue
		}
	}

	return failedImages
}
