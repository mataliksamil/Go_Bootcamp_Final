package product

import (
	pg "github.com/go-pg/pg/v10"
)

var dbConnect *pg.DB

func InitialProduct(db *pg.DB) {
	dbConnect = db
}
