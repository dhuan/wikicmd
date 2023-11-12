package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/dhuan/wikicmd/internal/config"
	"github.com/spf13/cobra"
)

var switchCmd = &cobra.Command{
	Use:   "switch",
	Short: "Switch between your available Wikis.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		beforeCommand(false)

		configRoot, err := config.GetAll()
		if err != nil {
			panic(err)
		}

		wikiId := args[0]

		wikiExists := false
		availableWikis := make([]string, len(configRoot.Wikis))
		for i, config := range configRoot.Wikis {
			if wikiId == config.Id {
				wikiExists = true
			}

			availableWikis[i] = config.Id
		}

		availableWikisFormatted := strings.Join(availableWikis, ",")

		if !wikiExists {
			fmt.Printf("No wiki exists with the given ID: %s\n\nThe available Wikis you can switch to are: %s\n", wikiId, availableWikisFormatted)

			os.Exit(1)
		}

		if configRoot.Default == wikiId {
			fmt.Println("Done!")

			os.Exit(0)
		}

		configRoot.Default = wikiId

		err = config.Set(configRoot)
		if err != nil {
			panic(err)
		}

		fmt.Println("Done!")
	},
}
