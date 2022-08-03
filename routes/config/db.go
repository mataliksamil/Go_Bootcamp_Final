package config

import (
	"log"
	"os"

	"github.com/go-pg/pg/v9"
)

// Connecting to db
func Connect() *pg.DB {
	opts := &pg.Options{
		User:     "postgres",
		Password: "123456",
		Addr:     "localhost:8080",
		Database: "gobootcamp",
	}
	var db *pg.DB = pg.Connect(opts)
	if db == nil {
		log.Printf("Failed to connect")
		os.Exit(100)
	}
	log.Printf("Connected to db")
	return db
}
