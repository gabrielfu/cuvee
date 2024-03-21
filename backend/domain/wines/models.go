package wines

type Vintage string

type ID string

type Purchase struct {
	ID       ID      `json:"id"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
	Date     string  `json:"date"`
}

type Wine struct {
	ID      ID         `json:"id"`
	Name    string     `json:"name"`
	Vintage Vintage    `json:"vintage"`
	Format  string     `json:"format"`
	Country string     `json:"country"`
	Region  string     `json:"region"`
	Details []Purchase `json:"details"` // purchase records
	Summary Purchase   `json:"summary"` // total quantity and average price
}
