package wines_test

import (
	"cuvee/domain/wines"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSummarizePurchases(t *testing.T) {
	purchases := []wines.Purchase{
		{Quantity: 1, Price: 15.0, Date: "2020-01-01"},
		{Quantity: 2, Price: 30.0, Date: "2020-01-02"},
		{Quantity: 3, Price: 25.0, Date: "2020-01-03"},
	}
	summary := wines.SummarizePurchases(purchases)
	assert.Equal(t, 6, summary.Quantity)
	assert.Equal(t, 25.0, summary.Price)
}
