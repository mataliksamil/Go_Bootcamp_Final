package services

import (
	"github.com/mataliksamil/Go_Bootcamp_Final/controllers"
)

var GIVENAMOUNT int = 500

type DiscountCalculator interface {
}

func Discount4thOrder(b []controllers.Basket) (bool, float64) {
	lastBasket := b[len(b)]
	discountOk := true
	basketCost := lastBasket.TotalCost
	for _, myb := range b {
		if myb.TotalCost > float64(GIVENAMOUNT) {
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

func Last3Order(u *controllers.User) {

}
