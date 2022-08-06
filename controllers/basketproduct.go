package controllers

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	pg "github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	guuid "github.com/google/uuid"
)

type BasketProduct struct {
	BasketProductID string `pg:",pk" json:"basketproduct_id"`
	ProductCount    int    `json:"product_count"`

	ProductID string   `json:"product_id"`
	Product   *Product `pg:"rel:has-one"`

	BasketID string  `json:"basket_id"`
	Basket   *Basket `pg:"rel:has-one"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// status eklenecek
	// cost
}

// Create BasketProduct Table
func CreateBasketProductTable(db *pg.DB) error {
	opts := &orm.CreateTableOptions{
		FKConstraints: true,
		IfNotExists:   true,
	}

	createError := db.Model(&BasketProduct{}).CreateTable(opts)
	if createError != nil {
		log.Printf("Error while creating basket table, Reason: %v\n", createError)
		return createError
	}
	log.Printf("Basket table created")
	return nil
}

func CreateBasketProduct(c *gin.Context) {
	var basketProduct BasketProduct
	c.BindJSON(&basketProduct)

	basketProduct_id := guuid.New().String() //
	product_id := basketProduct.ProductID    //
	basket_id := basketProduct.BasketID      //
	product_count := basketProduct.ProductCount

	_, insertError := dbConnect.Model(&BasketProduct{
		BasketProductID: basketProduct_id,
		ProductCount:    product_count,
		BasketID:        basket_id,
		ProductID:       product_id,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}).Insert()

	if insertError != nil {
		log.Printf("Error while inserting new Basket into db, Reason: %v\n", insertError)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Something went wrong",
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"message": "BasketProduct created Successfully",
	})

}

/*
func GetBasketProductsByBasketId(){

		basket_id := c.Param("basket_id")
		basket_product := &BasketProduct{Basket_ID: basket_id}

		err := dbConnect.Model(product).WherePK().Select()
		if err != nil {
			log.Printf("Error while getting a single todo, Reason: %v\n", err)
			c.JSON(http.StatusNotFound, gin.H{
				"status":  http.StatusNotFound,
				"message": "Todo not found",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "Single Product",
			"data":    product,
		})


} */
