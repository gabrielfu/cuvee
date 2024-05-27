package search

import (
	"context"
	"encoding/json"

	"google.golang.org/api/customsearch/v1"
	"google.golang.org/api/option"
)

type GoogleSearchEngine struct {
	cse *customsearch.CseService
	cx  string // Google Search cx
}

type GooglePageMap struct {
	Metatags []struct {
		OGDescription string `json:"og:description"`
	} `json:"metatags"`
}

func NewGoogleSearchEngine(cx string, apiKey string) (*GoogleSearchEngine, error) {
	s, err := customsearch.NewService(context.Background(), option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}
	cse := customsearch.NewCseService(s)
	return &GoogleSearchEngine{cse: cse, cx: cx}, nil
}

func (g *GoogleSearchEngine) doSearch(query string, searchType string) (*customsearch.Search, error) {
	return g.cse.List().Cx(g.cx).Q(query).SearchType(searchType).Do()
}

func (g *GoogleSearchEngine) WebSearch(query string, param any) (*WebSearchResult, error) {
	search, err := g.doSearch(query, "searchTypeUndefined")
	if err != nil {
		return nil, err
	}

	var items []WebSearchResultItem
	for i, item := range search.Items {
		var pagemap GooglePageMap
		if err := json.Unmarshal(item.Pagemap, &pagemap); err != nil {
			return nil, err
		}
		var description string
		if len(pagemap.Metatags) > 0 {
			description = pagemap.Metatags[0].OGDescription
		}
		items = append(items, WebSearchResultItem{
			Index:   i,
			Title:   item.Title,
			Link:    item.Link,
			Snippet: item.Snippet,
			Desc:    description,
		})
	}
	return &WebSearchResult{Items: items}, nil
}

func (g *GoogleSearchEngine) ImageSearch(query string, param any) (*ImageSearchResult, error) {
	search, err := g.doSearch(query, "image")
	if err != nil {
		return nil, err
	}

	var items []ImageSearchResultItem
	for i, item := range search.Items {
		items = append(items, ImageSearchResultItem{
			Index:  i,
			Link:   item.Link,
			Height: int(item.Image.Height),
			Width:  int(item.Image.Width),
		})
	}
	return &ImageSearchResult{Items: items}, nil
}
