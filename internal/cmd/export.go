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
		exportTypeValid := validateExportTypeFlag(flagExportType)
		if !exportTypeValid {
			fmt.Println(fmt.Sprintf("Type '%s' is not valid.\n\nThe valid types are: all,page,image.", flagExportType))

			os.Exit(1)
		}

		wikiConfig, apiCredentials, hook := beforeCommand()
		exportTo := args[0]

		exportCount := 0
		if flagExportType == export_type_all || flagExportType == export_type_page {
			exportCount, err = runExport(wikiConfig, apiCredentials, mw.NewStateForGetAllPages(), exportTo, 0, hook)
			if err != nil {
				panic(err)
			}
		}

		exportImagesCount := 0
		if flagExportType == export_type_all || flagExportType == export_type_image {
			exportImagesCount, err = runExportImages(wikiConfig, apiCredentials, exportTo, mw.FIRST_RUN, 0, hook)
			if err != nil {
				panic(err)
			}
		}

		fmt.Println(fmt.Sprintf("%d page(s) and %d and image(s) have been exported.\nDone!", exportCount, exportImagesCount))
	},
}

func runExportImages(
	config *mw.Config,
	apiCredentials *mw.ApiCredentials,
	exportTo string,
	continuation string,
	exportCount int,
	hook *mw.HookOptions,
) (int, error) {
	images, nextContinuation, finished, err := mw.GetAllImages(config, apiCredentials, continuation, hook)
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
	fmt.Println("Fetching next batch!!!!!!!!!!!!!!!!!!!!!")
	return runExportImages(config, apiCredentials, exportTo, nextContinuation, exportCount, hook)
}

func runExport(
	config *mw.Config,
	apiCredentials *mw.ApiCredentials,
	state map[string]string,
	exportTo string,
	exportCount int,
	hook *mw.HookOptions,
) (int, error) {
	pages, newState, finished, err := mw.GetAllPages(config, apiCredentials, state, hook)
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
	return runExport(config, apiCredentials, newState, exportTo, exportCount, hook)
}

func validateExportTypeFlag(exportType string) bool {
	return utils.AnyEquals(export_types, exportType)
}
