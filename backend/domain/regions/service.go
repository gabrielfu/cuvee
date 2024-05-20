package regions

import (
	"context"
)

type RegionService struct {
	regionRepo *RegionRepository
}

func NewRegionService(regionRepo *RegionRepository) *RegionService {
	return &RegionService{
		regionRepo: regionRepo,
	}
}

func (s *RegionService) ListRegions(ctx context.Context, wineID string) ([]Region, error) {
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

func (s *RegionService) GetRegion(ctx context.Context, wineID, symbol string) (*Region, error) {
	dao, err := s.regionRepo.GetRegion(ctx, wineID, symbol)
	if err != nil {
		return nil, err
	}
	region := daoToRegion(dao)
	return &region, nil
}

func (s *RegionService) CreateRegion(ctx context.Context, region *Region) error {
	return s.regionRepo.CreateRegion(ctx, regionToDAO(region))
}

func (s *RegionService) UpdateRegion(ctx context.Context, region *Region) error {
	return s.regionRepo.UpdateRegion(ctx, regionToDAO(region))
}

func (s *RegionService) DeleteRegion(ctx context.Context, wineID, symbol string) error {
	return s.regionRepo.DeleteRegion(ctx, wineID, symbol)
}
