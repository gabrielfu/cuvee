package wines

type PurchaseSummary struct {
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

type Purchase struct {
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
	Date     string  `json:"date"`
}

type Wine struct {
	ID      string          `json:"id"`
	Name    string          `json:"name"`
	Vintage string          `json:"vintage"`
	Format  string          `json:"format"`
	Country string          `json:"country"`
	Region  string          `json:"region"`
	Details []Purchase      `json:"details"` // purchase records
	Summary PurchaseSummary `json:"summary"` // total quantity and average price
}

func daoToWine(w WineDAO) Wine {
	var purchases []Purchase
	for i := range w.Purchases {
		purchases = append(purchases, daoToPurchase(w.Purchases[i]))
	}
	return Wine{
		ID:      w.ID.Hex(),
		Name:    w.Name,
		Vintage: w.Vintage,
		Format:  w.Format,
		Country: w.Country,
		Region:  w.Region,
		Details: purchases,
		Summary: SummarizePurchases(purchases),
	}
}

func wineToDAO(w Wine) WineDAO {
	var purchases []PurchaseDAO
	for i := range w.Details {
		purchases = append(purchases, purchaseToDAO(w.Details[i]))
	}
	return WineDAO{
		Name:      w.Name,
		Vintage:   w.Vintage,
		Format:    w.Format,
		Country:   w.Country,
		Region:    w.Region,
		Purchases: purchases,
	}
}

func daoToPurchase(p PurchaseDAO) Purchase {
	return Purchase{
		Quantity: p.Quantity,
		Price:    p.Price,
		Date:     p.Date,
	}
}

func purchaseToDAO(p Purchase) PurchaseDAO {
	return PurchaseDAO{
		Quantity: p.Quantity,
		Price:    p.Price,
		Date:     p.Date,
	}
}

func SummarizePurchases(purchases []Purchase) PurchaseSummary {
	var totalQuantity int
	var totalPrice float64
	for i := range purchases {
		totalQuantity += purchases[i].Quantity
		totalPrice += purchases[i].Price * float64(purchases[i].Quantity)
	}
	return PurchaseSummary{
		Quantity: totalQuantity,
		Price:    totalPrice / float64(totalQuantity),
	}
}
