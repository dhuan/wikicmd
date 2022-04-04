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
			handleErrorGettingApiCredentials(err, config.User, config.Address)

			os.Exit(1)
		}

		exportCount, err := runExport(&wikiConfig, apiCredentials, exportTo, mw.FIRST_RUN, 0)
		if err != nil {
			panic(err)
		}

		fmt.Println(fmt.Sprintf("%d page(s) have been exported.\nDone!", exportCount))
	},
}

func runExport(config *mw.Config, apiCredentials *mw.ApiCredentials, exportTo string, continuation string, exportCount int) (int, error) {
	pages, nextContinuation, finished, err := mw.GetAllPages(config, apiCredentials, continuation)
	exportCount = exportCount + len(pages)
	if err != nil {
		return 0, err
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
		return exportCount, nil
	}

	pages = make([]mw.Page, 0, 0)
	fmt.Println("Fetching next batch.")
	return runExport(config, apiCredentials, exportTo, nextContinuation, exportCount)
}
