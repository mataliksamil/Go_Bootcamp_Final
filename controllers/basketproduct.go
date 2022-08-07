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

func AddProductToBasket(c *gin.Context) {
	var basketProduct BasketProduct
	c.BindJSON(&basketProduct)

	basketProduct_id := guuid.New().String() //
	product_id := basketProduct.ProductID    //
	basket_id := basketProduct.BasketID      //
	product_count := basketProduct.ProductCount

	// update product stock accordingly
	err := ChangeProductStock(product_id, product_count)
	if err != nil {
		log.Printf("Add to basket failed, Reason: %v \n", err)
		return
	}
	// returns existing "Basket Product ID" and "product count" of it  if duplicate
	prevId, prevCount := AddIfDuplicate(basket_id, product_id)

	if prevId != "" { // update product count in PB if already this product can be found in the basket
		_, err := dbConnect.Model(&BasketProduct{}).Set("product_count = ?", product_count+prevCount).Where("basket_product_id = ?", prevId).Update()
		if err != nil {
			log.Printf("Error, Reason: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  500,
				"message": "Something went wrong",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status":  200,
			"message": "BasketProduct Edited Successfully",
		})

	} else { // create new basket product if there is no BP with the same product
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
}

// returns "Basket Product ID" if duplicate
func AddIfDuplicate(b_id string, p_id string) (string, int) {
	basket := &Basket{BasketID: b_id}
	err := dbConnect.Model(basket).Relation("BasketProducts").WherePK().Select()
	if err != nil {
		log.Printf("Error while getting the Basket, Reason: %v\n", err)
	} else {
		for _, bp := range basket.BasketProducts {
			if bp.ProductID == p_id {
				return bp.BasketProductID, bp.ProductCount
			}
		}
	}
	return "", 0
}
