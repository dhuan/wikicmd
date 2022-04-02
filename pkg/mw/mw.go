package mw

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"

	"github.com/dhuan/wikicmd/internal/utils"
)

type Config struct {
	BaseAddress string
	Login       string
	Password    string
}

type UploadWarning int

const (
	UPLOAD_WARNING_NONE                UploadWarning = iota
	UPLOAD_WARNING_SAME_FILE_NO_CHANGE               = iota
)

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

type UploadResult struct {
	Success bool
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

type uploadResponse struct {
	Upload struct {
		Result string `json:"result"`
	} `json:"upload"`
	Error struct {
		Code string `json:"code"`
		Info string `json:"info"`
	} `json:"error"`
}

type loginResponse struct {
	Login struct {
		Result     string `json:"result"`
		LgUserId   int    `json:"lguserid"`
		LgUsername string `json:"lgusername"`
	} `json:"login"`
}

var MAP_MW_ERROR_WARNING = map[UploadWarning]string{
	UPLOAD_WARNING_SAME_FILE_NO_CHANGE: "fileexists-no-change",
}

func GetApiCredentials(config *Config) (*ApiCredentials, error) {
	loginTokenSet, err := getLoginToken(config)
	if err != nil {
		return &ApiCredentials{}, errors.New("Failed to get login token.")
	}
	fmt.Println(fmt.Sprintf("Got Login Token Set\nCookie: %s\nToken:%s", loginTokenSet.Cookie, loginTokenSet.Token))

	loginResult, err := login(config, loginTokenSet)
	if err != nil {
		return &ApiCredentials{}, errors.New("Failed to login.")
	}
	fmt.Println(fmt.Sprintf("Got Login Result Set\nCookie: %s", loginResult.Cookie))

	csrfToken, err := getCsrfToken(config, loginTokenSet, loginResult)
	if err != nil {
		return &ApiCredentials{}, errors.New("Failed to get csrf token.")
	}
	fmt.Println(fmt.Sprintf("Got CSRF\nToken: %s", csrfToken.Token))

	return &ApiCredentials{CsrfToken: csrfToken, LoginResult: loginResult}, nil
}

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

func Upload(config *Config, credentials *ApiCredentials, fileName string, fileContent io.Reader) (error, []UploadWarning, bool) {
	buffer := &bytes.Buffer{}
	writer := multipart.NewWriter(buffer)

	values := map[string]string{
		"action":         "upload",
		"format":         "json",
		"filename":       fileName,
		"ignorewarnings": "1",
		"token":          credentials.CsrfToken.Token,
	}

	for key, value := range values {
		err := writer.WriteField(key, value)

		if err != nil {
			return err, []UploadWarning{}, false
		}
	}

	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return err, []UploadWarning{}, false
	}

	_, err = io.Copy(part, fileContent)
	if err != nil {
		return err, []UploadWarning{}, false
	}

	err = writer.Close()
	if err != nil {
		return err, []UploadWarning{}, false
	}

	request, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/api.php", config.BaseAddress),
		buffer,
	)
	if err != nil {
		return err, []UploadWarning{}, false
	}

	request.Header.Set("Content-Type", writer.FormDataContentType())
	request.Header.Set("Cookie", credentials.LoginResult.Cookie)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err, []UploadWarning{}, false
	}

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return err, []UploadWarning{}, false
	}

	decodedJson := &uploadResponse{}
	err = json.Unmarshal(bodyBytes, &decodedJson)
	if err != nil {
		return err, []UploadWarning{}, false
	}

	warning := resolveUploadWarningFromUploadResponse(decodedJson)
	if warning != UPLOAD_WARNING_NONE {
		return nil, []UploadWarning{warning}, false
	}

	if decodedJson.Error.Code != "" {
		return errors.New(fmt.Sprintf("%s: %s", decodedJson.Error.Code, decodedJson.Error.Info)), []UploadWarning{}, false
	}

	return nil, []UploadWarning{}, true
}

func resolveUploadWarningFromUploadResponse(response *uploadResponse) UploadWarning {
	if response.Error.Code == "" {
		return UPLOAD_WARNING_NONE
	}

	return utils.MapValueSearch[UploadWarning, string](MAP_MW_ERROR_WARNING, response.Error.Code, UPLOAD_WARNING_NONE)
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

func parseUploadResponse(decodedJson *uploadResponse, response *http.Response) (*UploadResult, error) {
	return &UploadResult{true}, nil
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
