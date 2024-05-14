package ratings

type Rating struct {
	Score    string
	Maturity string
	Notes    string
}

// a map of region to year to rating.
type RegionYearRatingMap map[string]map[string]Rating

type VintageChart interface {
	// ListRegions returns a list of regions.
	ListRegions() []string

	// GetRating returns a rating for the given region and year.
	// If the rating is not found, it returns an empty Rating object.
	GetRating(region string, year string) Rating
}
