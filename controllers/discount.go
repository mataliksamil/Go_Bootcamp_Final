package controllers

var GIVENAMOUNT float64 = 500
var MOUNTHLYTOTAL float64 = 2000

// Every fourth order whose total is more than given amount may have discount
func Discount4thOrder(user_id string) (float64, float64, error) {
	var myBasket = &Basket{}
	err := dbConnect.Model(myBasket).
		Relation("BasketProducts").
		Relation("BasketProducts.Product").
		Where("user_id=?", user_id).
		Where("basket_status=?", 1).
		Select()
	if err != nil {
		return 0.0, 0.0, err
	}
	basketCostVAT := myBasket.TotalVAT
	basketCost := myBasket.TotalCost
	// Iterate through baskets of last 3(passive) and current(active) basket

	lastBasketProducts := myBasket.BasketProducts
	for _, lbP := range lastBasketProducts {
		switch lbP.Product.VatRate {
		case 18:
			basketCostVAT = basketCostVAT - (lbP.Product.Price*float64(lbP.ProductCount))*0.15*1.18 // 15% discount for 18% vat
			basketCost = basketCost - (lbP.Product.Price*float64(lbP.ProductCount))*0.15            // 15% discount for 18% vat
		case 8:
			basketCostVAT = basketCostVAT - (lbP.Product.Price*float64(lbP.ProductCount))*0.10*1.08 // 15% discount for 18% vat
			basketCost = basketCost - (lbP.Product.Price*float64(lbP.ProductCount))*0.10            // 10% discount for 8% VAT
		}
	}
	return basketCost, basketCostVAT, nil
}

func Discount4thItem(user_id string) (float64, float64, error) {
	var myBasket = &Basket{}
	err := dbConnect.Model(myBasket).
		Relation("BasketProducts").
		Relation("BasketProducts.Product").
		Where("user_id=?", user_id).
		Where("basket_status=?", 1).
		Select()
	if err != nil {
		return 0.0, 0.0, err
	}
	basketProducts := myBasket.BasketProducts
	basketCost := myBasket.TotalCost
	basketCostVAT := myBasket.TotalVAT

	for _, myBP := range basketProducts {
		if myBP.ProductCount > 3 {
			basketCost = basketCost - float64(myBP.ProductCount-3)*0.08
			basketCostVAT = basketCostVAT - float64(myBP.ProductCount-3)*0.08
		}
	}
	return basketCostVAT, basketCostVAT, nil

}
