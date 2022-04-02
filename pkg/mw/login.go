package mw

import (
	"fmt"
	"net/url"
)

func getLoginToken(config *Config) (*LoginTokenSet, error) {
	return requestWrapper[loginTokenResponse, LoginTokenSet](
		fmt.Sprintf("%s/api.php?action=query&format=json&meta=tokens&type=login", config.BaseAddress),
		"GET",
		url.Values{},
		&loginTokenResponse{},
		&LoginTokenSet{},
		parseGetApiCredentials,
		map[string]string{},
	)
}

func login(config *Config, loginTokenSet *LoginTokenSet) (*LoginResult, error) {
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
	)
}

func getCsrfToken(config *Config, loginTokenSet *LoginTokenSet, loginResult *LoginResult) (*CsrfToken, error) {
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
	)
}
