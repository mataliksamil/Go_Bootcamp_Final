package basket

import (
	pg "github.com/go-pg/pg/v10"
)

var dbConnect *pg.DB

func InitialBasket(db *pg.DB) {
	dbConnect = db
}
