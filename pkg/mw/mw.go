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

type StateGetAllPages struct {
	FirstRun   bool
	Namespaces []string
	ApContinue string
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
	Exists  bool
}

type Image struct {
	Name    string
	Content []byte
}

type UploadResult struct {
	Success bool
}

type ApiCredentials struct {
	CsrfToken   *CsrfToken
	LoginResult *LoginResult
}

type HookOptions struct {
	BeforeRequest        func(string)
	AfterRequest         func(map[string]string)
	OnReceivedLoginToken func(map[string]string)
	OnLogin              func(map[string]string)
	OnCsrf               func(map[string]string)
}

var MAP_MW_ERROR_WARNING = map[UploadWarning]string{
	UPLOAD_WARNING_SAME_FILE_NO_CHANGE: "fileexists-no-change",
}

var FIRST_RUN = ""

func GetApiCredentials(config *Config, hook *HookOptions) (*ApiCredentials, error) {
	loginTokenSet, err := getLoginToken(config, hook)
	if err != nil {
		return &ApiCredentials{}, errors.New("Failed to get login token.")
	}
	hook.OnReceivedLoginToken(map[string]string{
		"cookie": loginTokenSet.Cookie,
		"token":  loginTokenSet.Token,
	})

	loginResult, err := login(config, loginTokenSet, hook)
	if err != nil {
		return &ApiCredentials{}, err
	}
	hook.OnLogin(map[string]string{
		"cookie": loginResult.Cookie,
	})

	csrfToken, err := getCsrfToken(config, loginTokenSet, loginResult, hook)
	if err != nil {
		return &ApiCredentials{}, errors.New("Failed to get csrf token.")
	}
	hook.OnCsrf(map[string]string{
		"token": csrfToken.Token,
	})

	return &ApiCredentials{CsrfToken: csrfToken, LoginResult: loginResult}, nil
}

func Edit(
	config *Config,
	credentials *ApiCredentials,
	title string,
	content string,
	summary string,
	hook *HookOptions,
) (*EditResult, error) {
	postFields := url.Values{
		"action":      {"edit"},
		"format":      {"jsonfm"},
		"title":       {title},
		"text":        {content},
		"wrappedhtml": {"1"},
		"token":       {credentials.CsrfToken.Token},
	}

	if summary != "" {
		postFields["summary"] = []string{summary}
	}

	return requestWrapper[editResponse, EditResult](
		fmt.Sprintf("%s/api.php", config.BaseAddress),
		"POST",
		postFields,
		&editResponse{},
		&EditResult{},
		parseEditResponse,
		map[string]string{
			"Cookie": credentials.LoginResult.Cookie,
		},
		hook,
	)
}

func GetAllImages(config *Config, credentials *ApiCredentials, continuation string, hook *HookOptions) ([]Image, string, bool, error) {
	images := make([]Image, 0)
	requestUrl := fmt.Sprintf("%s/api.php?action=query&format=json&list=allimages&ailimit=5", config.BaseAddress)

	if continuation != FIRST_RUN {
		requestUrl = fmt.Sprintf("%s&aicontinue=%s", requestUrl, continuation)
	}

	response, err := requestWrapper[getAllImagesResponse, getAllImagesResponse](
		requestUrl,
		"GET",
		url.Values{},
		&getAllImagesResponse{},
		&getAllImagesResponse{},
		func(decodedJson *getAllImagesResponse, response *http.Response) (*getAllImagesResponse, error) {
			return decodedJson, nil
		},
		map[string]string{},
		hook,
	)
	if err != nil {
		return []Image{}, continuation, true, err
	}

	for _, imageResponse := range response.Query.AllImages {
		imageContent, err := utils.Wget(imageResponse.Url)
		if err != nil {
			return []Image{}, "", true, err
		}

		images = append(images, Image{imageResponse.Name, imageContent})
	}

	finished := response.Continue.AiContinue == ""
	continuationNew := response.Continue.AiContinue

	return images, continuationNew, finished, nil
}

func parseStateGetAllPages(stateSerialized map[string]string) (*StateGetAllPages, error) {
	firstRun, ok := stateSerialized["firstRun"]
	if !ok {
		return &StateGetAllPages{}, fmt.Errorf("Field %s missing from serialized state.", "firstRun")
	}

	namespaces, ok := stateSerialized["namespaces"]
	if !ok {
		return &StateGetAllPages{}, fmt.Errorf("Field %s missing from serialized state.", "namespaces")
	}

	apContinue, ok := stateSerialized["apContinue"]
	if !ok {
		return &StateGetAllPages{}, fmt.Errorf("Field %s missing from serialized state.", "apContinue")
	}

	firstRunParsed := utils.StringToBool(firstRun)
	namespacesParsed := strings.Split(namespaces, ",")

	return &StateGetAllPages{firstRunParsed, namespacesParsed, apContinue}, nil
}

var (
	namespace_main      = "0"
	namespace_mediawiki = "8"
	namespace_template  = "10"
	namespace_category  = "14"
)

func NewStateForGetAllPages() map[string]string {
	return serializeStateForGetAllPages(&StateGetAllPages{
		FirstRun:   true,
		Namespaces: []string{namespace_main, namespace_mediawiki, namespace_template, namespace_category},
		ApContinue: "",
	})
}

func GetAllPages(
	config *Config,
	credentials *ApiCredentials,
	stateSerialized map[string]string,
	hook *HookOptions) ([]Page, map[string]string, bool, error) {
	requestUrl := fmt.Sprintf("%s/api.php?action=query&format=json&list=allpages&rawcontinue=1&aplimit=5", config.BaseAddress)

	state, err := parseStateGetAllPages(stateSerialized)
	if err != nil {
		return []Page{}, map[string]string{}, true, err
	}

	apNamespace := state.Namespaces[0]
	if !state.FirstRun {
		requestUrl = fmt.Sprintf("%s&apcontinue=%s&apnamespace=%s", requestUrl, state.ApContinue, apNamespace)
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
		hook,
	)
	if err != nil {
		return []Page{}, stateSerialized, true, err
	}

	pages := make([]Page, 0)
	for _, page := range response.Query.AllPages {
		fetchedPage, err := GetPage(config, credentials, utils.FormatPageNameInput(page.Title), hook)

		if err != nil {
			return []Page{}, stateSerialized, true, err
		}

		pages = append(pages, Page{page.Title, fetchedPage.Content, true})
	}

	finished, newState, err := getNewStateForGetAllPages(response, state)
	if err != nil {
		return []Page{}, stateSerialized, true, err
	}

	return pages, serializeStateForGetAllPages(newState), finished, nil
}

func serializeStateForGetAllPages(state *StateGetAllPages) map[string]string {
	return map[string]string{
		"firstRun":   utils.BoolToString(state.FirstRun),
		"namespaces": strings.Join(state.Namespaces, ","),
		"apContinue": state.ApContinue,
	}
}

func getNewStateForGetAllPages(response *getAllPagesResponse, state *StateGetAllPages) (bool, *StateGetAllPages, error) {
	apContinue := response.QueryContinue.AllPages.ApContinue
	finishedCurrent := apContinue == ""
	finishedAll := finishedCurrent && len(state.Namespaces) == 1

	if finishedAll {
		return true, &StateGetAllPages{
			FirstRun:   false,
			Namespaces: []string{},
			ApContinue: "",
		}, nil
	}

	if !finishedCurrent {
		return false, &StateGetAllPages{
			FirstRun:   false,
			Namespaces: state.Namespaces,
			ApContinue: apContinue,
		}, nil
	}

	return false, &StateGetAllPages{
		FirstRun:   false,
		Namespaces: utils.RemoveNth[string](state.Namespaces, 0),
		ApContinue: apContinue,
	}, nil
}

func GetPage(config *Config, credentials *ApiCredentials, title string, hook *HookOptions) (*Page, error) {
	return requestWrapper[getPageResponse, Page](
		fmt.Sprintf("%s/api.php?action=parse&format=json&page=%s&prop=wikitext&formatversion=2", config.BaseAddress, title),
		"GET",
		url.Values{},
		&getPageResponse{},
		&Page{},
		parseGetPageResponse(title),
		map[string]string{},
		hook,
	)
}

func Upload(config *Config, credentials *ApiCredentials, fileName string, fileContent io.Reader) ([]UploadWarning, bool, error) {
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
			return []UploadWarning{}, false, err
		}
	}

	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return []UploadWarning{}, false, err
	}

	_, err = io.Copy(part, fileContent)
	if err != nil {
		return []UploadWarning{}, false, err
	}

	err = writer.Close()
	if err != nil {
		return []UploadWarning{}, false, err
	}

	request, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/api.php", config.BaseAddress),
		buffer,
	)
	if err != nil {
		return []UploadWarning{}, false, err
	}

	request.Header.Set("Content-Type", writer.FormDataContentType())
	request.Header.Set("Cookie", credentials.LoginResult.Cookie)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return []UploadWarning{}, false, err
	}

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return []UploadWarning{}, false, err
	}

	decodedJson := &uploadResponse{}
	err = json.Unmarshal(bodyBytes, &decodedJson)
	if err != nil {
		return []UploadWarning{}, false, err
	}

	warning := resolveUploadWarningFromUploadResponse(decodedJson)
	if warning != UPLOAD_WARNING_NONE {
		return []UploadWarning{warning}, false, nil
	}

	if decodedJson.Error.Code != "" {
		return []UploadWarning{}, false, fmt.Errorf("%s: %s", decodedJson.Error.Code, decodedJson.Error.Info)
	}

	return []UploadWarning{}, true, nil
}

func resolveUploadWarningFromUploadResponse(response *uploadResponse) UploadWarning {
	if response.Error.Code == "" {
		return UPLOAD_WARNING_NONE
	}

	return utils.MapValueSearch[UploadWarning, string](MAP_MW_ERROR_WARNING, response.Error.Code, UPLOAD_WARNING_NONE)
}
