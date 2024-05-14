package search

import "net/http"

type SearchEngine interface {
	Search(query string, param any) (*http.Response, error)
}
