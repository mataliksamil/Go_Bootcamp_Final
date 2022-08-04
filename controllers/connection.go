package controllers

import pg "github.com/go-pg/pg/v9"

var dbConnect *pg.DB

func InitiateDB(db *pg.DB) {
	dbConnect = db
}
