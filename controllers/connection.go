package controllers

import pg "github.com/go-pg/pg/v10"

var dbConnect *pg.DB

func InitiateDB(db *pg.DB) {
	dbConnect = db
}
