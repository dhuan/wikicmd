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

func beforeCommand() (*mw.Config, *mw.ApiCredentials, *mw.HookOptions) {
	config, err := config.Get()
	if err != nil {
		panic(err)
	}

	hookBeforeRequest := nilHook
	hookOnReceivedLoginToken := nilHookWithParams
	hookOnLogin := nilHookWithParams
	hookOnCsrf := nilHookWithParams
	if flagVerbose {
		hookBeforeRequest = logHook("Requesting %s")
		hookOnReceivedLoginToken = logWithParamsHook("Got login token set.\n* Cookie: %s\n* Token %s", []string{"cookie", "token"})
		hookOnLogin = logWithParamsHook("Got Login Result Set\n* Cookie: %s", []string{"cookie"})
		hookOnCsrf = logWithParamsHook("Got CSRF Token: %s", []string{"token"})
	}

	hook := &mw.HookOptions{
		BeforeRequest:        hookBeforeRequest,
		OnReceivedLoginToken: hookOnReceivedLoginToken,
		OnLogin:              hookOnLogin,
		OnCsrf:               hookOnCsrf,
	}

	wikiConfig := &mw.Config{
		BaseAddress: config.Address,
		Login:       config.User,
		Password:    config.Password,
	}

	apiCredentials, err := mw.GetApiCredentials(wikiConfig, hook)
	if err != nil {
		handleErrorGettingApiCredentials(err, config.User, config.Address)

		os.Exit(1)
	}

	return wikiConfig, apiCredentials, hook
}

func nilHook(logMessage string) {
}

func nilHookWithParams(params map[string]string) {
}

func logHook(message string) func(string) {
	return func(logMessage string) {
		fmt.Println(fmt.Sprintf(message, logMessage))
	}
}

func logWithParamsHook(message string, paramNames []string) func(map[string]string) {
	return func(params map[string]string) {
		paramValues := make([]interface{}, len(paramNames))
		for i, paramName := range paramNames {
			paramValue, ok := params[paramName]
			if !ok {
				paramValues[i] = "unknown"

				continue
			}

			paramValues[i] = paramValue
		}

		fmt.Println(fmt.Sprintf(message, paramValues...))
	}
}
