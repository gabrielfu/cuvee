package vintagecharts

import (
	"encoding/json"
	"os"
	"reflect"
	"regexp"
	"sort"
)

type RPProvider struct {
	ratings         RegionYearRatingMap // Ratings is a map of region to year to rating.
	regions         []string            // Regions is a list of regions.
	maturityLegends map[string]string
}

type RPRating struct {
	Group       string `json:"group"`
	SubGroup    string `json:"subGroup"`
	Country     string `json:"country"`
	Region      string `json:"region"`
	RegionLabel string `json:"regionLabel"`
	Year        string `json:"year"`
	Rating      string `json:"rating"`
}

func readChartFile(chartFile string) ([]RPRating, error) {
	f, err := os.Open(chartFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var ratings []RPRating
	if err := json.NewDecoder(f).Decode(&ratings); err != nil {
		return nil, err
	}
	return ratings, nil
}

// sortedMapKeys returns a sorted list of keys from a map.
func sortedMapKeys(m interface{}) (keyList []string) {
	keys := reflect.ValueOf(m).MapKeys()
	for _, key := range keys {
		keyList = append(keyList, key.Interface().(string))
	}
	sort.Strings(keyList)
	return
}

func loadChartFile(chartFile string) (RegionYearRatingMap, error) {
	rpRatings, err := readChartFile(chartFile)
	if err != nil {
		return nil, err
	}

	ratingRexgep := regexp.MustCompile(`([0-9\- ]+)?(C|E|NV|I|NT|R|T)?`)

	ratings := make(RegionYearRatingMap)
	for _, r := range rpRatings {
		if _, ok := ratings[r.RegionLabel]; !ok {
			ratings[r.RegionLabel] = make(map[string]Rating)
		}

		matches := ratingRexgep.FindStringSubmatch(r.Rating)
		score := matches[1]
		maturity := matches[2]

		ratings[r.RegionLabel][r.Year] = Rating{
			Region:   r.RegionLabel,
			Year:     r.Year,
			Score:    score,
			Maturity: maturity,
			Notes:    "",
		}
	}
	return ratings, err
}

func loadMaturityLegends(maturityFile string) (map[string]string, error) {
	f, err := os.Open(maturityFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var legends map[string]string
	if err := json.NewDecoder(f).Decode(&legends); err != nil {
		return nil, err
	}
	return legends, nil
}

func NewRPProvider(chartFile string, maturityFile string) (*RPProvider, error) {
	ratings, err := loadChartFile(chartFile)
	if err != nil {
		return nil, err
	}

	regions := sortedMapKeys(ratings)

	// Load the maturity legends.
	maturityLegends, err := loadMaturityLegends(maturityFile)
	if err != nil {
		return nil, err
	}

	return &RPProvider{
		ratings:         ratings,
		regions:         regions,
		maturityLegends: maturityLegends,
	}, nil
}

func (r *RPProvider) Name() string {
	return "Robert Parker"
}

func (r *RPProvider) Symbol() string {
	return "RP"
}

func (r *RPProvider) ListRegions() []string {
	return r.regions
}

func (r *RPProvider) GetRating(region string, year string) Rating {
	rating, ok := r.ratings[region][year]
	if !ok {
		return Rating{}
	}
	// translate maturity
	rating.Maturity = r.maturityLegends[rating.Maturity]
	return rating
}
