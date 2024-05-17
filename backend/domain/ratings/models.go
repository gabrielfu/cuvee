package ratings

type VintageChart struct {
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
}

type Region struct {
	ID       string `json:"id,omitempty"`
	WineID   string `json:"wineId"`
	VCSymbol string `json:"vcSymbol"`
	Region   string `json:"region"`
}

func regionToDAO(region *Region) RegionDAO {
	return RegionDAO{
		WineID:   region.WineID,
		VCSymbol: region.VCSymbol,
		Region:   region.Region,
	}
}

func daoToRegion(region *RegionDAO) Region {
	return Region{
		ID:       region.ID.Hex(),
		WineID:   region.WineID,
		VCSymbol: region.VCSymbol,
		Region:   region.Region,
	}
}

type SuggestRegionRequest struct {
	WineID   string `json:"wineId"`
	VCSymbol string `json:"vcSymbol"`
}
