package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/dhuan/wikicmd/internal/config"
	"github.com/dhuan/wikicmd/pkg/mw"
)

func handleErrorGettingApiCredentials(err error, user string, wikiAddress string) {
	if errors.Is(err, mw.ErrLogin) {
		fmt.Println(fmt.Sprintf("Failed to login with user %s in %s.", user, wikiAddress))

		return
	}

	fmt.Println("Failed to get API Credentials.")
}

func beforeCommand() (*mw.Config, *mw.ApiCredentials) {
	config, err := config.Get()
	if err != nil {
		panic(err)
	}

	wikiConfig := &mw.Config{
		BaseAddress: config.Address,
		Login:       config.User,
		Password:    config.Password,
	}

	apiCredentials, err := mw.GetApiCredentials(wikiConfig)
	if err != nil {
		handleErrorGettingApiCredentials(err, config.User, config.Address)

		os.Exit(1)
	}

	return wikiConfig, apiCredentials
}
