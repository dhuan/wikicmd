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
	Name    string
	Content string
}

type UploadResult struct {
	Success bool
}

type ApiCredentials struct {
	CsrfToken   *CsrfToken
	LoginResult *LoginResult
}

var MAP_MW_ERROR_WARNING = map[UploadWarning]string{
	UPLOAD_WARNING_SAME_FILE_NO_CHANGE: "fileexists-no-change",
}

var FIRST_RUN = ""

func GetApiCredentials(config *Config) (*ApiCredentials, error) {
	loginTokenSet, err := getLoginToken(config)
	if err != nil {
		return &ApiCredentials{}, errors.New("Failed to get login token.")
	}
	fmt.Println(fmt.Sprintf("Got Login Token Set\nCookie: %s\nToken:%s", loginTokenSet.Cookie, loginTokenSet.Token))

	loginResult, err := login(config, loginTokenSet)
	if err != nil {
		return &ApiCredentials{}, err
	}
	fmt.Println(fmt.Sprintf("Got Login Result Set\nCookie: %s", loginResult.Cookie))

	csrfToken, err := getCsrfToken(config, loginTokenSet, loginResult)
	if err != nil {
		return &ApiCredentials{}, errors.New("Failed to get csrf token.")
	}
	fmt.Println(fmt.Sprintf("Got CSRF\nToken: %s", csrfToken.Token))

	return &ApiCredentials{CsrfToken: csrfToken, LoginResult: loginResult}, nil
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

func GetAllPages(config *Config, credentials *ApiCredentials, continuation string) ([]Page, string, bool, error) {
	requestUrl := fmt.Sprintf("%s/api.php?action=query&format=json&list=allpages&rawcontinue=1&aplimit=5", config.BaseAddress)

	if continuation != FIRST_RUN {
		requestUrl = fmt.Sprintf("%s&apcontinue=%s", requestUrl, continuation)
	}

	response, err := requestWrapper[getAllPagesResponse, getAllPagesResponse](
		requestUrl,
		"GET",
		url.Values{},
		&getAllPagesResponse{},
		&getAllPagesResponse{},
		func(decodedJson *getAllPagesResponse, response *http.Response) (*getAllPagesResponse, error) {
			return decodedJson, nil
		},
		map[string]string{},
	)
	if err != nil {
		return []Page{}, continuation, true, err
	}

	pages := make([]Page, 0)
	for _, page := range response.Query.AllPages {
		fetchedPage, err := GetPage(config, credentials, utils.FormatPageNameInput(page.Title))

		if err != nil {
			return []Page{}, continuation, true, err
		}

		pages = append(pages, Page{page.Title, fetchedPage.Content})
	}

	finished := response.QueryContinue.AllPages.ApContinue == ""
	continuationNew := response.QueryContinue.AllPages.ApContinue

	return pages, continuationNew, finished, nil
}

func GetPage(config *Config, credentials *ApiCredentials, title string) (*Page, error) {
	return requestWrapper[getPageResponse, Page](
		fmt.Sprintf("%s/api.php?action=parse&format=json&page=%s&prop=wikitext&formatversion=2", config.BaseAddress, title),
		"GET",
		url.Values{},
		&getPageResponse{},
		&Page{},
		parseGetPageResponse(title),
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
