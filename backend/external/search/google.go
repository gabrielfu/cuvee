package search

import (
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

const baseUrl = "https://www.googleapis.com/customsearch/v1"

func NewGoogleSearchEngine(googleSearchCx string, googleSearchApiKey string) *GoogleSearchEngine {
	return &GoogleSearchEngine{googleSearchCx: googleSearchCx, googleSearchApiKey: googleSearchApiKey}
}

func (g GoogleSearchEngine) Search(query string, param any) (*http.Response, error) {
	u, _ := url.Parse(baseUrl)

	q := url.Values{}
	q.Set("q", url.QueryEscape(query))
	q.Set("cx", g.googleSearchCx)
	q.Set("key", g.googleSearchApiKey)
	if searchType := param.(GoogleSearchParam).searchType; searchType != "" {
		q.Set("searchType", searchType)
	}

	u.RawQuery = q.Encode()
	return http.Get(u.String())
}
