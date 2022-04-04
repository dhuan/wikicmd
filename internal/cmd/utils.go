package cmd

import (
	"errors"
	"fmt"

	"github.com/dhuan/wikicmd/pkg/mw"
)

func handleErrorGettingApiCredentials(err error, user string, wikiAddress string) {
	if errors.Is(err, mw.ErrLogin) {
		fmt.Println(fmt.Sprintf("Failed to login with user %s in %s.", user, wikiAddress))

		return
	}

	fmt.Println("Failed to get API Credentials.")
}
