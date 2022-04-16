package mw

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func requestWrapper[D interface{}, T interface{}](
	url string,
	method string,
	postData url.Values,
	obj *D,
	result *T,
	parse func(obj *D, response *http.Response) (*T, error),
	headers map[string]string,
	hook *HookOptions,
) (*T, error) {
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

	hook.BeforeRequest(buildLogMessageBeforeRequest(method, url, postData))

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

func buildLogMessageBeforeRequest(method, url string, postData url.Values) string {
	if len(postData) == 0 {
		return fmt.Sprintf("%s %s", method, url)
	}

	return fmt.Sprintf("%s %s\n%s", method, url, urlValuesToString(postData))
}

func urlValuesToString(urlValues url.Values) string {
	stringList := make([]string, len(urlValues))

	i := 0
	for key, values := range urlValues {
		stringList[i] = fmt.Sprintf("* %s: %s", key, strings.Join(values, " "))

		i++
	}

	return strings.Join(stringList, "\n")
}
