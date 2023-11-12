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
		fmt.Printf("Failed to login with user %s in %s.\n", user, wikiAddress)

		return
	}

	fmt.Println("Failed to get API Credentials.")
}

func beforeCommand(withApiCredentials bool) (*mw.Config, *mw.ApiCredentials, *config.UserSettings, *mw.HookOptions) {
	userConfig, configRoot, err := config.Get()
	if errors.Is(err, config.ErrConfigDoesNotExist) {
		fmt.Println("You don't seem to have a configuration file. Try 'wikicmd config' to initialize a new configuration.")

		os.Exit(1)
	}
	if errors.Is(err, config.ErrConfigDoesNotHaveWiki) {
		fmt.Println("Your configuration doesn't seem to have a Wiki defined. Try 'wikicmd config --new' to initialize a new configuration.")

		os.Exit(1)
	}
	if err != nil {
		panic(err)
	}

	hookBeforeRequest := nilHook
	hookAfterRequest := nilHookWithParams
	hookOnReceivedLoginToken := nilHookWithParams
	hookOnLogin := nilHookWithParams
	hookOnCsrf := nilHookWithParams
	if flagVerbose {
		hookBeforeRequest = logHook("Requesting %s")
		hookAfterRequest = logWithParamsHook("Request finished.\n* %s", []string{"responseBody"})
		hookOnReceivedLoginToken = logWithParamsHook("Got login token set.\n* Cookie: %s\n* Token %s", []string{"cookie", "token"})
		hookOnLogin = logWithParamsHook("Got Login Result Set\n* Cookie: %s", []string{"cookie"})
		hookOnCsrf = logWithParamsHook("Got CSRF Token: %s", []string{"token"})
	}

	hook := &mw.HookOptions{
		BeforeRequest:        hookBeforeRequest,
		AfterRequest:         hookAfterRequest,
		OnReceivedLoginToken: hookOnReceivedLoginToken,
		OnLogin:              hookOnLogin,
		OnCsrf:               hookOnCsrf,
	}

	wikiConfig := &mw.Config{
		BaseAddress: userConfig.Address,
		Login:       userConfig.User,
		Password:    userConfig.Password,
	}

	apiCredentials := &mw.ApiCredentials{}
	if withApiCredentials {
		apiCredentials, err = mw.GetApiCredentials(wikiConfig, hook)
		if err != nil {
			handleErrorGettingApiCredentials(err, userConfig.User, userConfig.Address)

			os.Exit(1)
		}
	}

	userSettings := config.GetUserSettings(configRoot)

	return wikiConfig, apiCredentials, userSettings, hook
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
