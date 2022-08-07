package services

import (
	"github.com/mataliksamil/Go_Bootcamp_Final/controllers"
)

var GIVENAMOUNT float64 = 500
var MOUNTHLYTOTAL float64 = 2000

type DiscountCalculator interface {
}

// Every fourth order whose total is more than given amount may have discount
func Discount4thOrder(b []controllers.Basket) (bool, float64) {
	lastBasket := b[len(b)]
	discountOk := true
	basketCost := lastBasket.TotalCost

	// Iterate through baskets of last 3(passive) and current(active) basket
	for _, myb := range b {
		if myb.TotalCost > GIVENAMOUNT {
			continue
		} else {
			return false, lastBasket.TotalCost
		}
	}
	lastBasketProducts := lastBasket.BasketProducts
	for _, lbP := range lastBasketProducts {
		switch lbP.Product.VatRate {
		case 18:
			basketCost = basketCost - (lbP.Product.Price*float64(lbP.ProductCount))*0.15
		case 8:
			basketCost = basketCost - (lbP.Product.Price*float64(lbP.ProductCount))*0.10
		}
	}
	return discountOk, basketCost
}

func Discount4thItem(b controllers.Basket) (bool, float64) {
	basketProducts := b.BasketProducts
	basketCost := b.TotalCost
	discountOk := false
	for _, myBP := range basketProducts {
		if myBP.ProductCount > 3 {
			basketCost = basketCost - (float64(myBP.ProductCount)-3)*0.08
			discountOk = true
		}
	}
	return discountOk, basketCost

}

func DiscountByMonthly(b []controllers.Basket) (bool, float64) {
	lastBasket := b[len(b)]
	monthlyTotal := 0.0

	for _, myB := range b {
		monthlyTotal = monthlyTotal + myB.TotalCost
	}

	if monthlyTotal >= MOUNTHLYTOTAL { // discount applied
		return true, lastBasket.TotalCost * 0.9
	}

	return false, lastBasket.TotalCost //	discount does not applied

}
