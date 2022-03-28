package cmd

import (
	"fmt"

	"github.com/dhuan/wikicmd/internal/config"
	"github.com/dhuan/wikicmd/pkg/editor"
	"github.com/dhuan/wikicmd/pkg/mw"
	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit pages",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pageName := args[0]

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

		page, err := mw.GetPage(&wikiConfig, &mw.ApiCredentials{CsrfToken: csrfToken, LoginResult: loginResult}, pageName)
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
			&wikiConfig,
			&mw.ApiCredentials{CsrfToken: csrfToken, LoginResult: loginResult},
			pageName,
			newContent,
		)
		if err != nil {
			panic("Failed to edit.")
		}
	},
}
