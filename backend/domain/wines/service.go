package wines

import "context"

type WineService struct {
	repo     *WineRepository
	validate *WineJSONValidator
}

func NewWineService(repo *WineRepository, validate *WineJSONValidator) *WineService {
	return &WineService{repo: repo, validate: validate}
}

func withDefaultWineImage(wine Wine) Wine {
	if wine.ImageUrl == "" {
		wine.ImageUrl = "https://static.vecteezy.com/system/resources/previews/016/475/672/original/transparent-dark-wine-bottle-with-blank-label-and-burgundy-foil-capsule-seal-png.png"
	}
	return wine
}

func (s *WineService) CreateWine(ctx context.Context, wine *Wine) (string, error) {
	dao := wineToDAO(withDefaultWineImage(*wine))
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
	wines := make([]Wine, len(daos))
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
