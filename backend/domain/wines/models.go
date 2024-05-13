package wines

type PurchaseSummary struct {
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

type Purchase struct {
	Quantity int     `json:"quantity" validate:"gte=1"`
	Price    float64 `json:"price" validate:"gte=0"`
	Date     string  `json:"date" validate:"required"`
}

type Wine struct {
	ID        string          `json:"id"`
	Name      string          `json:"name" validate:"required"`
	Vintage   string          `json:"vintage" validate:"required,vintage"`
	Format    string          `json:"format" validate:"required"`
	Country   string          `json:"country" validate:"required"`
	Region    string          `json:"region" validate:"required"`
	Purchases []Purchase      `json:"purchases" validate:"dive,required"` // purchase records
	Summary   PurchaseSummary `json:"summary"`                            // total quantity and average price
	ImageUrl  string          `json:"imageUrl"`
}

func daoToWine(w WineDAO) Wine {
	var purchases []Purchase
	for i := range w.Purchases {
		purchases = append(purchases, daoToPurchase(w.Purchases[i]))
	}
	return Wine{
		ID:        w.ID.Hex(),
		Name:      w.Name,
		Vintage:   w.Vintage,
		Format:    w.Format,
		Country:   w.Country,
		Region:    w.Region,
		Purchases: purchases,
		Summary:   SummarizePurchases(purchases),
		ImageUrl:  w.ImageUrl,
	}
}

func wineToDAO(w Wine) WineDAO {
	var purchases []PurchaseDAO
	for i := range w.Purchases {
		purchases = append(purchases, purchaseToDAO(w.Purchases[i]))
	}
	return WineDAO{
		Name:      w.Name,
		Vintage:   w.Vintage,
		Format:    w.Format,
		Country:   w.Country,
		Region:    w.Region,
		Purchases: purchases,
		ImageUrl:  w.ImageUrl,
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
