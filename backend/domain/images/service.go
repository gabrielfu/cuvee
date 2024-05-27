package images

import (
	"context"
	"cuvee/external/search"
	"net/http"
	"strings"
	"sync"
)

type ImageService struct {
	searchEngine search.SearchEngine
}

type ImageSearchRequest struct {
	Name    string `json:"name"`
	Vintage string `json:"vintage,omitempty"`
	Country string `json:"country,omitempty"`
	Region  string `json:"region,omitempty"`
}

func NewImageService(searchEngine search.SearchEngine) *ImageService {
	return &ImageService{searchEngine: searchEngine}
}

func validateIsImage(link string) bool {
	resp, err := http.Get(link)
	if err != nil || resp.StatusCode != http.StatusOK {
		return false
	}
	contentType := resp.Header.Get("Content-Type")
	return strings.HasPrefix(contentType, "image/")
}

func validateImageLink(link string) bool {
	return strings.HasPrefix(link, "https://") && validateIsImage(link)
}

func buildQuery(request ImageSearchRequest) string {
	query := "vivino "
	if request.Name != "" {
		query += request.Name + " "
	}
	if request.Vintage != "" {
		query += request.Vintage + " "
	}
	if request.Country != "" {
		query += request.Country + " "
	}
	if request.Region != "" {
		query += request.Region + " "
	}
	return query
}

func (s *ImageService) Search(ctx context.Context, request ImageSearchRequest) (*search.ImageSearchResult, error) {
	query := buildQuery(request)
	searchResponse, err := s.searchEngine.ImageSearch(query, nil)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	valid := make([]bool, len(searchResponse.Items))
	for i, respItem := range searchResponse.Items {
		wg.Add(1)
		_i := i
		go func(id int, item search.ImageSearchResultItem) {
			defer wg.Done()
			if validateImageLink(item.Link) {
				valid[id] = true
			}
		}(_i, respItem)
	}
	wg.Wait()

	var filtered []search.ImageSearchResultItem
	for i, item := range searchResponse.Items {
		if valid[i] {
			filtered = append(filtered, item)
		}
	}
	searchResponse.Items = filtered
	return searchResponse, nil
}
