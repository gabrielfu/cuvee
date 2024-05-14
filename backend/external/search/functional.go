package search

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Search is a generic function that takes a SearchEngine and a query string
// and unmarshals the JSON response into the provided type.
func Search[T any](s SearchEngine, query string, param any) (*T, error) {
	resp, err := s.Search(query, param)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("received %d response code from Google Image Search: %s", resp.StatusCode, string(bodyBytes))
	}

	var searchResponse T
	if err := json.NewDecoder(resp.Body).Decode(&searchResponse); err != nil {
		return nil, err
	}
	return &searchResponse, nil
}
