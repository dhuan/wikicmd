package mw

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Config struct {
	BaseAddress string
	Login       string
	Password    string
}

type LoginTokenSet struct {
	Token  string
	Cookie string
}

type LoginResult struct {
	Cookie string
}

type CsrfToken struct {
	Token string
}

type EditResult struct {
	Success bool
}

type Page struct {
	Content string
}

type ApiCredentials struct {
	CsrfToken   *CsrfToken
	LoginResult *LoginResult
}

type loginTokenResponse struct {
	Query struct {
		Tokens struct {
			Logintoken string `json:"logintoken"`
		} `json:"tokens"`
	} `json:"query"`
}

type csrfTokenResponse struct {
	Query struct {
		Tokens struct {
			CsrfToken string `json:"csrftoken"`
		} `json:"tokens"`
	} `json:"query"`
}

type editResponse struct {
	Status int `json:"status"`
}

type getPageResponse struct {
	Parse struct {
		Wikitext string `json:"wikitext"`
	} `json:"parse"`
}

type loginResponse struct {
	Login struct {
		Result     string `json:"result"`
		LgUserId   int    `json:"lguserid"`
		LgUsername string `json:"lgusername"`
	} `json:"login"`
}

func GetLoginToken(config *Config) (*LoginTokenSet, error) {
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

func Login(config *Config, loginTokenSet *LoginTokenSet) (*LoginResult, error) {
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

func GetCsrfToken(config *Config, loginTokenSet *LoginTokenSet, loginResult *LoginResult) (*CsrfToken, error) {
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

func Edit(config *Config, credentials *ApiCredentials, title string, content string) (*EditResult, error) {
	return requestWrapper[editResponse, EditResult](
		fmt.Sprintf("%s/api.php", config.BaseAddress),
		"POST",
		url.Values{
			"action":      {"edit"},
			"format":      {"jsonfm"},
			"title":       {title},
			"text":        {content},
			"summary":     {"test summary"},
			"wrappedhtml": {"1"},
			"token":       {credentials.CsrfToken.Token},
		},
		&editResponse{},
		&EditResult{},
		parseEditResponse,
		map[string]string{
			"Cookie": credentials.LoginResult.Cookie,
		},
	)
}

func GetPage(config *Config, credentials *ApiCredentials, title string) (*Page, error) {
	return requestWrapper[getPageResponse, Page](
		fmt.Sprintf("%s/api.php?action=parse&format=json&page=%s&prop=wikitext&formatversion=2", config.BaseAddress, title),
		"GET",
		url.Values{},
		&getPageResponse{},
		&Page{},
		parseGetPageResponse,
		map[string]string{},
	)
}

func parseGetApiCredentials(decodedJson *loginTokenResponse, response *http.Response) (*LoginTokenSet, error) {
	token := decodedJson.Query.Tokens.Logintoken
	cookie := response.Header.Get("Set-Cookie")

	return &LoginTokenSet{token, cookie}, nil
}

func parseLoginResponse(decodedJson *loginResponse, response *http.Response) (*LoginResult, error) {
	cookie := response.Header.Get("Set-Cookie")

	return &LoginResult{cookie}, nil
}

func parseCsrfTokenResponse(decodedJson *csrfTokenResponse, response *http.Response) (*CsrfToken, error) {
	return &CsrfToken{decodedJson.Query.Tokens.CsrfToken}, nil
}

func parseEditResponse(decodedJson *editResponse, response *http.Response) (*EditResult, error) {
	return &EditResult{true}, nil
}

func parseGetPageResponse(decodedJson *getPageResponse, response *http.Response) (*Page, error) {
	return &Page{decodedJson.Parse.Wikitext}, nil
}

func requestWrapper[D interface{}, T interface{}](
	url string,
	method string,
	postData url.Values,
	obj *D,
	result *T,
	parse func(obj *D, response *http.Response,
	) (*T, error), headers map[string]string) (*T, error) {
	var response *http.Response
	var err error

	client := &http.Client{}

	req, err := http.NewRequest(method, url, strings.NewReader(postData.Encode()))
	if err != nil {
		return result, err
	}

	if method == "POST" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	for headerKey, headerValue := range headers {
		req.Header.Set(headerKey, headerValue)
	}

	response, err = client.Do(req)
	if err != nil {
		return result, err
	}

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return result, err
	}

	decodedJson := obj

	err = json.Unmarshal(bodyBytes, &decodedJson)
	if err != nil {
		return result, err
	}

	return parse(obj, response)
}
