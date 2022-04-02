package mw

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
