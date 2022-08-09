package config

import (
	pg "github.com/go-pg/pg/v10"
	basket "github.com/mataliksamil/Go_Bootcamp_Final/controllers/basket"
	basketProduct "github.com/mataliksamil/Go_Bootcamp_Final/controllers/basketProduct"
	product "github.com/mataliksamil/Go_Bootcamp_Final/controllers/product"
	user "github.com/mataliksamil/Go_Bootcamp_Final/controllers/user"
)

func InitiateDB(db *pg.DB) {
	basket.InitialBasket(db)
	basketProduct.InitialBasketProduct(db)
	product.InitialProduct(db)
	user.InitialUser(db)
}
