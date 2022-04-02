package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dhuan/wikicmd/internal/config"
	"github.com/dhuan/wikicmd/internal/utils"
	"github.com/dhuan/wikicmd/pkg/mw"
	"github.com/spf13/cobra"
)

var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload images",
	Run: func(cmd *cobra.Command, filePaths []string) {
		config, err := config.Get()
		if err != nil {
			panic(err)
		}

		wikiConfig := mw.Config{
			BaseAddress: config.Address,
			Login:       config.User,
			Password:    config.Password,
		}

		loginTokenSet, err := mw.GetLoginToken(&wikiConfig)
		if err != nil {
			panic("Failed to get login token.")
		}
		fmt.Println(fmt.Sprintf("Got Login Token Set\nCookie: %s\nToken:%s", loginTokenSet.Cookie, loginTokenSet.Token))

		loginResult, err := mw.Login(&wikiConfig, loginTokenSet)
		if err != nil {
			panic("Failed to login.")
		}
		fmt.Println(fmt.Sprintf("Got Login Result Set\nCookie: %s", loginResult.Cookie))

		csrfToken, err := mw.GetCsrfToken(&wikiConfig, loginTokenSet, loginResult)
		if err != nil {
			panic("Failed to get csrf token.")
		}
		fmt.Println(fmt.Sprintf("Got CSRF\nToken: %s", csrfToken.Token))

		failedImages := validateImages(filePaths)

		if len(failedImages) > 0 {
			failedImagesStr := strings.Join(failedImages, "\n")
			panic(fmt.Sprintf("The following files cannot be uploaded:\n%s", failedImagesStr))
		}

		for _, filePath := range filePaths {
			fmt.Println(fmt.Sprintf("Uploading %s", filePath))

			fileContent, err := os.Open(filePath)
			if err != nil {
				panic(fmt.Sprintf("Failed to upload: %s", filePath))
			}

			fileName := filepath.Base(filePath)

			err = mw.Upload(&wikiConfig, &mw.ApiCredentials{CsrfToken: csrfToken, LoginResult: loginResult}, fileName, fileContent)
			if err != nil {
				panic(err)
			}
		}
	},
}

func validateImages(filePaths []string) []string {
	failedImages := make([]string, 0)

	for _, filePath := range filePaths {
		if extensionIsValid := utils.ExtensionMatches([]string{
			"png",
			"jpg",
			"jpeg",
			"gif",
		}, filePath); !extensionIsValid {
			failedImages = append(failedImages, filePath)

			continue
		}
	}

	return failedImages
}
