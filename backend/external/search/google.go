package search

import (
	"fmt"
	"net/http"
	"net/url"
)

type GoogleSearchEngine struct {
	googleSearchCx     string
	googleSearchApiKey string
}

type GoogleSearchParam struct {
	searchType string
}

var ImageSearchGoogleSearchParam = GoogleSearchParam{searchType: "image"}
var WebSearchGoogleSearchParam = GoogleSearchParam{searchType: ""}

func paramToQuery(param GoogleSearchParam) string {
	if param.searchType == "" {
		return ""
	}
	return fmt.Sprintf("searchType=%s", param.searchType)
}

func NewGoogleSearchEngine(googleSearchCx string, googleSearchApiKey string) *GoogleSearchEngine {
	return &GoogleSearchEngine{googleSearchCx: googleSearchCx, googleSearchApiKey: googleSearchApiKey}
}

func (g GoogleSearchEngine) Search(query string, param any) (*http.Response, error) {
	reqUrl := fmt.Sprintf("https://www.googleapis.com/customsearch/v1?q=%s&cx=%s&key=%s", url.QueryEscape(query), g.googleSearchCx, g.googleSearchApiKey)
	paramQuery := paramToQuery(param.(GoogleSearchParam))
	if paramQuery != "" {
		reqUrl += "&" + paramQuery
	}
	return http.Get(reqUrl)
}
