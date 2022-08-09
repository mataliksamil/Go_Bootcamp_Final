package config

import (
	"log"
	"os"

	pg "github.com/go-pg/pg/v10"
	basket "github.com/mataliksamil/Go_Bootcamp_Final/controllers/basket"
	basketProduct "github.com/mataliksamil/Go_Bootcamp_Final/controllers/basketProduct"
	product "github.com/mataliksamil/Go_Bootcamp_Final/controllers/product"
	user "github.com/mataliksamil/Go_Bootcamp_Final/controllers/user"
)

// Connecting to db
func Connect() *pg.DB {
	opts := &pg.Options{
		User:     "postgres",
		Password: "123456",
		Addr:     "localhost:5432",
		Database: "go_database",
	}
	var db *pg.DB = pg.Connect(opts)
	if db == nil {
		log.Printf("Failed to connect")
		os.Exit(100)
	}
	log.Printf("Connected to db")

	/* 	err := controllers.CreateSchema(db)
	   	if err != nil {
	   		panic(err)
	   	}
	   	log.Printf("Schema created") */

	//controllers.CreateUserTable(db)
	user.CreateUserTable(db)
	basket.CreateBasketTable(db)
	product.CreateProductTable(db)
	basketProduct.CreateBasketProductTable(db)

	return db
}
