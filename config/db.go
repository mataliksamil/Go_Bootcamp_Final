package config

import (
	"log"
	"os"

	"github.com/go-pg/pg/v9"
	controllers "github.com/mataliksamil/Go_Bootcamp_Final/controllers"
	//	products "github.com/mataliksamil/Go_Bootcamp_Final/products"
	//users "github.com/mataliksamil/Go_Bootcamp_Final/users"
)

// Connecting to db
func Connect() *pg.DB {
	opts := &pg.Options{
		User:     "postgres",
		Password: "123456",
		Addr:     "localhost:5432",
		Database: "gobootcamp",
	}
	var db *pg.DB = pg.Connect(opts)
	if db == nil {
		log.Printf("Failed to connect")
		os.Exit(100)
	}
	log.Printf("Connected to db")
	controllers.CreateUserTable(db)
	//users.CreateUserTable(db)
	//	products.CreateProductTable(db)
	controllers.InitiateDB(db)

	return db
}
