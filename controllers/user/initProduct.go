package user

import (
	pg "github.com/go-pg/pg/v10"
)

var dbConnect *pg.DB

func InitialUser(db *pg.DB) {
	dbConnect = db
}
