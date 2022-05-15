package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/dhuan/wikicmd/internal/config"
	"github.com/dhuan/wikicmd/internal/input"
	"github.com/dhuan/wikicmd/pkg/editor"
	"github.com/spf13/cobra"
)

var defaultWiki = "https://en.wikipedia.org/w"

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Edit your configuration or create a new one",
	Run: func(cmd *cobra.Command, filePaths []string) {
		configFilePath, configFileExists, err := config.GetConfigFilePath()
		if err != nil {
			panic(err)
		}

		if !configFileExists || flagConfigNew {
			confirmed, err := newConfigWizard(configFilePath, flagConfigNew)
			if err != nil {
				panic(err)
			}

			if !confirmed {
				fmt.Println("Aborted!")
				os.Exit(1)
			}
		}

		userSettings := config.DefaultUserSettings
		configRoot := &config.ConfigRoot{}
		if configFileExists {
			configRoot, err = config.GetConfigFromPath(configFilePath)
		}
		if err != nil {
			panic(err)
		}

		userSettings = config.GetUserSettings(configRoot)

		err = editor.EditFile(userSettings.Editor, configFilePath)
		if err != nil {
			panic(err)
		}

		if configFileExists {
			fmt.Println("Done!")
		} else {
			fmt.Println(fmt.Sprintf("Done!\n\nRun \"wikicmd config\" again whenever you want to edit that configuration file."))
		}
	},
}

func newConfigWizard(filePath string, requestingNew bool) (bool, error) {
	if requestingNew {
		fmt.Println("Let's create a new configuration file. Be aware that your existing configuration file will be overwritten at the end of this process.")
	} else {
		fmt.Println("You don't seem to have any configuration file. Let's create one.")
	}

	fmt.Println("")
	inputWikiAddress := input.TextPrompt(fmt.Sprintf("Wiki address: (%s) ", defaultWiki), defaultWiki)
	inputLogin := input.TextPrompt("Login: ", "")
	inputPassword := input.TextPrompt("Password: ", "")

	fmt.Println(fmt.Sprintf("\nNext, a configuration file will be created for you and saved as %s\n", filePath))

	inputConfirm, _ := input.ResponsePrompt[bool](
		"Is this OK? (yes): ",
		map[string]bool{"yes": true},
		false,
		true,
	)

	if !inputConfirm {
		return false, nil
	}

	newConfig := config.ConfigRoot{
		Config: []config.Config{
			config.Config{
				Id:               "my_wiki",
				Address:          inputWikiAddress,
				User:             inputLogin,
				Password:         inputPassword,
				ImportExtensions: []string{},
			},
		},
		Default: "my_wiki",
	}

	jsonEncoded, err := json.MarshalIndent(newConfig, "", "  ")
	if err != nil {
		return false, err
	}

	os.WriteFile(filePath, jsonEncoded, 0644)

	return true, nil
}
