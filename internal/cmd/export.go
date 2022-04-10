package cmd

import (
	"fmt"
	"os"

	"github.com/dhuan/wikicmd/internal/utils"
	"github.com/dhuan/wikicmd/pkg/mw"
	"github.com/spf13/cobra"
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export pages",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		exportTypeValid := validateExportTypeFlag(FlagExportType)
		if !exportTypeValid {
			fmt.Println(fmt.Sprintf("Type '%s' is not valid.\n\nThe valid types are: all,page,image.", FlagExportType))

			os.Exit(1)
		}

		wikiConfig, apiCredentials := beforeCommand()
		exportTo := args[0]

		exportCount := 0
		if FlagExportType == export_type_all || FlagExportType == export_type_page {
			exportCount, err = runExport(wikiConfig, apiCredentials, exportTo, mw.FIRST_RUN, 0)
			if err != nil {
				panic(err)
			}
		}

		exportImagesCount := 0
		if FlagExportType == export_type_all || FlagExportType == export_type_image {
			exportImagesCount, err = runExportImages(wikiConfig, apiCredentials, exportTo, mw.FIRST_RUN, 0)
			if err != nil {
				panic(err)
			}
		}

		fmt.Println(fmt.Sprintf("%d page(s) and %d and image(s) have been exported.\nDone!", exportCount, exportImagesCount))
	},
}

func runExportImages(config *mw.Config, apiCredentials *mw.ApiCredentials, exportTo string, continuation string, exportCount int) (int, error) {
	images, nextContinuation, finished, err := mw.GetAllImages(config, apiCredentials, continuation)
	exportCount = exportCount + len(images)
	if err != nil {
		return 0, err
	}

	for _, image := range images {
		saveAs := fmt.Sprintf("%s/%s", exportTo, image.Name)
		fmt.Println(fmt.Sprintf("Writing %s", saveAs))
		if err := os.WriteFile(
			saveAs,
			image.Content,
			0644,
		); err != nil {
			panic(err)
		}
	}

	if finished {
		return exportCount, nil
	}

	images = make([]mw.Image, 0, 0)
	fmt.Println("Fetching next batch.")
	return runExportImages(config, apiCredentials, exportTo, nextContinuation, exportCount)
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

func validateExportTypeFlag(exportType string) bool {
	return utils.AnyEquals(export_types, exportType)
}
