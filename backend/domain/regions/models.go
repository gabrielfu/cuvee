package regions

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
