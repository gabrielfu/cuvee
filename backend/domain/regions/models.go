package regions

type Region struct {
	ID     string `json:"id,omitempty"`
	WineID string `json:"wineId"`
	Symbol string `json:"symbol"`
	Region string `json:"region"`
}

func regionToDAO(region *Region) RegionDAO {
	return RegionDAO{
		WineID: region.WineID,
		Symbol: region.Symbol,
		Region: region.Region,
	}
}

func daoToRegion(region *RegionDAO) Region {
	return Region{
		ID:     region.ID.Hex(),
		WineID: region.WineID,
		Symbol: region.Symbol,
		Region: region.Region,
	}
}
