package search

type ImageSearchResultItem struct {
	Index  int    `json:"index"`
	Link   string `json:"link"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
}

type ImageSearchResult struct {
	Items []ImageSearchResultItem `json:"items"`
}

type WebSearchResultItem struct {
	Index   int    `json:"index"`
	Title   string `json:"title"`
	Link    string `json:"link"`
	Snippet string `json:"snippet"`
	Desc    string `json:"desc"` // description field in meta tag or open graph tag
}

type WebSearchResult struct {
	Items []WebSearchResultItem `json:"items"`
}

type SearchEngine interface {
	WebSearch(query string, param any) (*WebSearchResult, error)
	ImageSearch(query string, param any) (*ImageSearchResult, error)
}
