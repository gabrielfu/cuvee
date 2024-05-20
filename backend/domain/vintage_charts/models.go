package vintagecharts

type VintageChart struct {
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
}

type Rating struct {
	Region   string
	Year     string
	Score    string
	Maturity string
	Notes    string
}
