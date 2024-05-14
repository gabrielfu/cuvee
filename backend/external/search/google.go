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

type GoogleImageSearchResponseItem struct {
	Link  string `json:"link"`
	Image struct {
		ContextLink string `json:"contextLink"`
		Height      int    `json:"height"`
		Width       int    `json:"width"`
	} `json:"image"`
}

type GoogleImageSearchResponse struct {
	Items []GoogleImageSearchResponseItem `json:"items"`
}

type GoogleWebSearchResponseItem struct {
	Title   string `json:"title"`
	Link    string `json:"link"`
	Snippet string `json:"snippet"`
	PageMap struct {
		Metatags []struct {
			OGDescription string `json:"og:description"`
		} `json:"metatags"`
	} `json:"pagemap"`
}

type GoogleWebSearchResponse struct {
	Items []GoogleWebSearchResponseItem `json:"items"`
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
