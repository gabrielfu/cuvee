package images

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
)

type ImageService struct {
	googleSearchCx     string
	googleSearchApiKey string
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

func NewImageService(googleSearchCx string, googleSearchApiKey string) *ImageService {
	return &ImageService{googleSearchCx: googleSearchCx, googleSearchApiKey: googleSearchApiKey}
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
	return url.QueryEscape(query)
}

func (s *ImageService) Search(ctx context.Context, request ImageSearchRequest) (ImageSearchResponse, error) {
	query := buildQuery(request)
	reqUrl := "https://content-customsearch.googleapis.com/customsearch/v1?searchType=image&q=" + query + "&cx=" + s.googleSearchCx + "&key=" + s.googleSearchApiKey
	resp, err := http.Get(reqUrl)

	if resp.StatusCode != http.StatusOK {
		p := make([]byte, resp.ContentLength)
		resp.Body.Read(p)
		msg := fmt.Errorf("received %d response code from Google Image Search: %s", resp.StatusCode, string(p))
		log.Printf("Google Image Search error (status=%d): %v\n", resp.StatusCode, msg)
		return ImageSearchResponse{}, err
	}

	if err != nil {
		return ImageSearchResponse{}, err
	}
	defer resp.Body.Close()

	var searchResponse GoogleImageSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResponse); err != nil {
		log.Printf("JSON decode Error: %v\n", err)
		return ImageSearchResponse{}, err
	}

	output := ImageSearchResponse{}
	output.Items = make([]ImageSearchResponseItem, 0)
	ch := make(chan ImageSearchResponseItem)
	var wg sync.WaitGroup // Add wait group

	for i, respItem := range searchResponse.Items {
		wg.Add(1) // Increment wait group counter
		go func(id int, item GoogleImageSearchResponseItem) {
			defer wg.Done() // Decrement wait group counter
			if !validateImageLink(item.Link) {
				return
			}
			ch <- ImageSearchResponseItem{
				ID:     id,
				Link:   item.Link,
				Height: item.Image.Height,
				Width:  item.Image.Width,
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
