package ratings

import (
	"context"
	"cuvee/domain/wines"
	"cuvee/external/llm"
	"cuvee/external/search"
	"fmt"
)

type RatingService struct {
	vcs        []VintageChartData
	llm        llm.LLM
	search     search.SearchEngine
	regionRepo *RegionRepository
}

func NewRatingService(vcs []VintageChartData, llm llm.LLM, search search.SearchEngine, regionRepo *RegionRepository) *RatingService {
	return &RatingService{
		vcs:        vcs,
		llm:        llm,
		search:     search,
		regionRepo: regionRepo,
	}
}

func (s *RatingService) ListRegions(ctx context.Context, wineID string) ([]Region, error) {
	daos, err := s.regionRepo.ListRegions(ctx, wineID)
	if err != nil {
		return nil, err
	}
	regions := make([]Region, 0)
	for i := range daos {
		regions = append(regions, daoToRegion(&daos[i]))
	}
	return regions, nil
}

func (s *RatingService) GetRegion(ctx context.Context, wineID, vcSymbol string) (*Region, error) {
	dao, err := s.regionRepo.GetRegion(ctx, wineID, vcSymbol)
	if err != nil {
		return nil, err
	}
	region := daoToRegion(dao)
	return &region, nil
}

func (s *RatingService) CreateRegion(ctx context.Context, region *Region) error {
	return s.regionRepo.CreateRegion(ctx, regionToDAO(region))
}

func (s *RatingService) UpdateRegion(ctx context.Context, region *Region) error {
	return s.regionRepo.UpdateRegion(ctx, regionToDAO(region))
}

func (s *RatingService) DeleteRegion(ctx context.Context, wineID, vcSymbol string) error {
	return s.regionRepo.DeleteRegion(ctx, wineID, vcSymbol)
}

func (s *RatingService) ListVintageCharts() []VintageChart {
	vcs := make([]VintageChart, 0)
	for _, vc := range s.vcs {
		vcs = append(vcs, VintageChart{
			Name:   vc.Name(),
			Symbol: vc.Symbol(),
		})
	}
	return vcs
}

func (s *RatingService) getVintageChartData(vcSymbol string) (VintageChartData, error) {
	for _, vc := range s.vcs {
		if vc.Symbol() == vcSymbol {
			return vc, nil
		}
	}
	return nil, fmt.Errorf("vintage chart not found: %s", vcSymbol)
}

func (s *RatingService) SuggestRegion(ctx context.Context, wine *wines.Wine, vcSymbol string) (*Region, error) {
	vc, err := s.getVintageChartData(vcSymbol)
	if err != nil {
		return nil, err
	}
	region, err := PickRegion(
		wine.Name,
		wine.Vintage,
		wine.Country,
		wine.Region,
		vc.ListRegions(),
		s.llm,
		s.search,
	)
	if err != nil {
		return nil, err
	}
	return &Region{
		WineID:   wine.ID,
		VCSymbol: vcSymbol,
		Region:   region,
	}, nil
}
