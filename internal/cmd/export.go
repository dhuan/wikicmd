package cmd

import (
	"fmt"
	"os"

	"github.com/dhuan/wikicmd/internal/config"
	"github.com/dhuan/wikicmd/pkg/mw"
	"github.com/spf13/cobra"
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export pages",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		exportTo := args[0]

		config, err := config.Get()
		if err != nil {
			panic(err)
		}

		wikiConfig := mw.Config{
			BaseAddress: config.Address,
			Login:       config.User,
			Password:    config.Password,
		}

		apiCredentials, err := mw.GetApiCredentials(&wikiConfig)
		if err != nil {
			panic(err)
		}

		err = runExport(&wikiConfig, apiCredentials, exportTo, mw.FIRST_RUN)
		if err != nil {
			panic(err)
		}
	},
}

func runExport(config *mw.Config, apiCredentials *mw.ApiCredentials, exportTo string, continuation string) error {
	pages, nextContinuation, finished, err := mw.GetAllPages(config, apiCredentials, continuation)
	if err != nil {
		return err
	}

	for _, page := range pages {
		fileName := fmt.Sprintf("%s/%s.wikitext", exportTo, page.Name)
		fmt.Println(fmt.Sprintf("Writing %s", fileName))
		if err := os.WriteFile(
			fileName,
			[]byte(page.Content),
			0644,
		); err != nil {
			panic(err)
		}
	}

	if finished {
		return nil
	}

	pages = make([]mw.Page, 0, 0)

	return runExport(config, apiCredentials, exportTo, nextContinuation)
}
