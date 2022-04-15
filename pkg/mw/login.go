package mw

import (
	"fmt"
	"net/url"
)

func getLoginToken(config *Config, hook *HookOptions) (*LoginTokenSet, error) {
	return requestWrapper[loginTokenResponse, LoginTokenSet](
		fmt.Sprintf("%s/api.php?action=query&format=json&meta=tokens&type=login", config.BaseAddress),
		"GET",
		url.Values{},
		&loginTokenResponse{},
		&LoginTokenSet{},
		parseGetApiCredentials,
		map[string]string{},
		hook,
	)
}

func login(config *Config, loginTokenSet *LoginTokenSet, hook *HookOptions) (*LoginResult, error) {
	return requestWrapper[loginResponse, LoginResult](
		fmt.Sprintf("%s/api.php", config.BaseAddress),
		"POST",
		url.Values{
			"format":     {"json"},
			"action":     {"login"},
			"lgname":     {config.Login},
			"lgpassword": {config.Password},
			"lgtoken":    {loginTokenSet.Token},
		},
		&loginResponse{},
		&LoginResult{},
		parseLoginResponse,
		map[string]string{
			"Cookie": loginTokenSet.Cookie,
		},
		hook,
	)
}

func getCsrfToken(config *Config, loginTokenSet *LoginTokenSet, loginResult *LoginResult, hook *HookOptions) (*CsrfToken, error) {
	return requestWrapper[csrfTokenResponse, CsrfToken](
		fmt.Sprintf("%s/api.php?action=query&format=json&meta=tokens", config.BaseAddress),
		"GET",
		url.Values{},
		&csrfTokenResponse{},
		&CsrfToken{},
		parseCsrfTokenResponse,
		map[string]string{
			"Cookie": loginResult.Cookie,
		},
		hook,
	)
}
