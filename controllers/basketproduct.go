package controllers

import (
	"fmt"
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
	var err error
	// update product stock accordingly

	if err != nil {
		errString := fmt.Sprintf("Add to basket not allowed, Reason : %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusMethodNotAllowed,
			"message": errString,
		})
		return
	}
	// returns existing "Basket Product ID" and "product count" of it  if duplicate
	prevId, prevCount := AddIfDuplicate(basket_id, product_id)

	if prevId != "" { // update product count in Product Basket(PB) if already this product can be found in the basket
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

func DiscardBasketProduct(c *gin.Context) {
	basketproduct_id := c.Param("basketproduct_id") // takes bp_id parameter
	var basketProduct BasketProduct
	c.BindJSON(&basketProduct)                   // json parser
	discardedCount := basketProduct.ProductCount // product count will be deleted, from the json

	prevCount, productID := PrevBPCount(basketproduct_id)
	if prevCount < discardedCount {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  400,
			"message": "There is not enough product in this basket",
		})
		return
	}
	_ = ChangeProductStock(productID, -discardedCount)
	if prevCount == discardedCount {
		DeleteBP(basketproduct_id, c)
	} else {
		UpdateBPCount(basketproduct_id, prevCount-discardedCount, c)
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

func PrevBPCount(bp_id string) (int, string) { // returns the basketProducts current product count
	basketProduct := &BasketProduct{BasketProductID: bp_id}
	err := dbConnect.Model(basketProduct).WherePK().Select()
	if err != nil {
		log.Printf("Error while getting the Basket Product , Reason: %v\n", err)
	}
	return basketProduct.ProductCount, basketProduct.ProductID
}

func DeleteBP(bp_id string, c *gin.Context) {
	basketproduct_id := c.Param("basketproduct_id") // takes bp_id parameter
	myBasketProduct := &BasketProduct{BasketProductID: basketproduct_id}

	_, err := dbConnect.Model(myBasketProduct).WherePK().Delete()
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
		"message": "BasketProduct Deleted Successfully",
	})
}

func UpdateBPCount(bp_id string, product_count int, c *gin.Context) {
	basketproduct_id := c.Param("basketproduct_id") // takes bp_id parameter
	myBasketProduct := &BasketProduct{BasketProductID: basketproduct_id}

	_, err := dbConnect.Model(myBasketProduct).Set("product_count = ?", product_count).WherePK().Update()
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
}
