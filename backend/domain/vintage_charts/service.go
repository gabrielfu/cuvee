package vintagecharts

import (
	"context"
	"cuvee/external/llm"
	"cuvee/external/search"
	"fmt"
)

type VintageChartService struct {
	providers []Provider
	llm       llm.LLM
	search    search.SearchEngine
}

type SuggestRequest struct {
	Name    string `json:"name"`
	Vintage string `json:"vintage,omitempty"`
	Country string `json:"country,omitempty"`
	Region  string `json:"region,omitempty"`
}

func NewVintageChartService(providers []Provider, llm llm.LLM, search search.SearchEngine) *VintageChartService {
	return &VintageChartService{
		providers: providers,
		llm:       llm,
		search:    search,
	}
}

func (s *VintageChartService) ListVintageCharts() []VintageChart {
	vcs := make([]VintageChart, 0)
	for _, provider := range s.providers {
		vcs = append(vcs, providerToVintageChart(provider))
	}
	return vcs
}

func (s *VintageChartService) getProvider(symbol string) (Provider, error) {
	for _, vc := range s.providers {
		if vc.Symbol() == symbol {
			return vc, nil
		}
	}
	return nil, fmt.Errorf("vintage chart not found: %s", symbol)
}

func (s *VintageChartService) ListRegions(symbol string) ([]string, error) {
	vc, err := s.getProvider(symbol)
	if err != nil {
		return nil, err
	}
	return vc.ListRegions(), nil
}

func (s *VintageChartService) SuggestRegion(ctx context.Context, request SuggestRequest, regions []string) (string, error) {
	region, err := PickRegion(request, regions, s.llm, s.search)
	if err != nil {
		return "", err
	}
	return region, nil
}

func (s *VintageChartService) GetRating(ctx context.Context, symbol, region, vintage string) (*Rating, error) {
	provider, err := s.getProvider(symbol)
	if err != nil {
		return nil, err
	}
	rating := provider.GetRating(region, vintage)
	return &rating, nil
}
