package ratings

import (
	"context"
	"cuvee/external/llm"
	"cuvee/external/search"
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
