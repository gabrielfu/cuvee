package wines

import "context"

type WineService struct {
	repo *WineRepository
}

func NewWineService(repo *WineRepository) *WineService {
	return &WineService{repo}
}

func (s *WineService) CreateWine(ctx context.Context, wine *Wine) (string, error) {
	dao := wineToDAO(*wine)
	return s.repo.Create(ctx, &dao)
}

func (s *WineService) GetWine(ctx context.Context, id string) (Wine, error) {
	dao, err := s.repo.Get(ctx, id)
	if err != nil {
		return Wine{}, err
	}
	return daoToWine(dao), nil
}

func (s *WineService) ListWines(ctx context.Context) ([]Wine, error) {
	daos, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	var wines []Wine
	for i := range daos {
		wines = append(wines, daoToWine(daos[i]))
	}
	return wines, nil
}

func (s *WineService) UpdateWine(ctx context.Context, id string, wine *Wine) error {
	dao := wineToDAO(*wine)
	return s.repo.Update(ctx, id, &dao)
}

func (s *WineService) DeleteWine(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
