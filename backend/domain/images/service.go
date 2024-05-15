package images

import (
	"context"
	"cuvee/external/search"
	"net/http"
	"sort"
	"strings"
	"sync"
)

type ImageService struct {
	searchEngine search.SearchEngine
}

type ImageSearchRequest struct {
	name    string
	vintage string
	country string
	region  string
}

type ImageSearchResponseItem struct {
	ID     int    `json:"id"`
	Link   string `json:"link"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
}

type ImageSearchResponse struct {
	Items []ImageSearchResponseItem `json:"items"`
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
	if request.name != "" {
		query += request.name + " "
	}
	if request.vintage != "" {
		query += request.vintage + " "
	}
	if request.country != "" {
		query += request.country + " "
	}
	if request.region != "" {
		query += request.region + " "
	}
	return query
}

func (s *ImageService) Search(ctx context.Context, request ImageSearchRequest) (ImageSearchResponse, error) {
	query := buildQuery(request)
	searchResponse, err := s.searchEngine.ImageSearch(query, nil)
	if err != nil {
		return ImageSearchResponse{}, err
	}

	output := ImageSearchResponse{Items: make([]ImageSearchResponseItem, 0)}
	ch := make(chan ImageSearchResponseItem)
	var wg sync.WaitGroup // Add wait group

	for i, respItem := range searchResponse.Items {
		wg.Add(1) // Increment wait group counter
		go func(id int, item search.ImageSearchResultItem) {
			defer wg.Done() // Decrement wait group counter
			if !validateImageLink(item.Link) {
				return
			}
			ch <- ImageSearchResponseItem{
				ID:     id,
				Link:   item.Link,
				Height: item.Height,
				Width:  item.Width,
			}
		}(i, respItem)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for item := range ch {
		output.Items = append(output.Items, item)
	}
	sort.Slice(output.Items, func(i, j int) bool {
		return output.Items[i].ID < output.Items[j].ID
	})
	return output, nil
}
