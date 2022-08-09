package basketProduct

import (
	pg "github.com/go-pg/pg/v10"
)

var dbConnect *pg.DB

func InitialBasketProduct(db *pg.DB) {
	dbConnect = db
}
