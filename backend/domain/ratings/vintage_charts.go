package ratings

type Rating struct {
	Region   string
	Year     string
	Score    string
	Maturity string
	Notes    string
}

// a map of region to year to rating.
type RegionYearRatingMap map[string]map[string]Rating

type VintageChart interface {
	Name() string

	Symbol() string

	// ListRegions returns a list of regions.
	ListRegions() []string

	// GetRating returns a rating for the given region and year.
	// If the rating is not found, it returns an empty Rating object.
	GetRating(region string, year string) Rating
}
