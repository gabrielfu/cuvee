package vintagecharts

// a map of region to vintage to rating.
type RegionVintageRatingMap map[string]map[string]Rating

type Provider interface {
	Name() string

	Symbol() string

	// ListRegions returns a list of regions.
	ListRegions() []string

	// GetRating returns a rating for the given region and vintage.
	// If the rating is not found, it returns an empty Rating object.
	GetRating(region string, vintage string) Rating
}

func providerToVintageChart(provider Provider) VintageChart {
	return VintageChart{
		Name:   provider.Name(),
		Symbol: provider.Symbol(),
	}
}
