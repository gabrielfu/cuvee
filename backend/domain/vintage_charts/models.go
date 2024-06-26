package vintagecharts

type VintageChart struct {
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
}

type Rating struct {
	Region   string `json:"region"`
	Vintage  string `json:"vintage"`
	Score    string `json:"score"`
	Maturity string `json:"maturity"`
	Notes    string `json:"notes"`
}
