package mw

import "net/http"

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
