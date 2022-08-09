package basketProduct

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	pg "github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	guuid "github.com/google/uuid"
	entities "github.com/mataliksamil/Go_Bootcamp_Final/entities"
)

// Create BasketProduct Table
func CreateBasketProductTable(db *pg.DB) error {
	opts := &orm.CreateTableOptions{
		FKConstraints: true,
		IfNotExists:   true,
	}

	createError := db.Model(&entities.BasketProduct{}).CreateTable(opts)
	if createError != nil {
		log.Printf("Error while creating basket table, Reason: %v\n", createError)
		return createError
	}
	log.Printf("Basket table created")
	return nil
}

func AddProductToBasket(c *gin.Context) {
	var basketProduct entities.BasketProduct
	c.BindJSON(&basketProduct)

	basketProduct_id := guuid.New().String() //
	product_id := basketProduct.ProductID    //
	basket_id := basketProduct.BasketID      //
	product_count := basketProduct.ProductCount

	// update product stock accordingly

	// returns existing "Basket Product ID" and "product count" of it  if duplicate
	prevId, prevCount := AddIfDuplicate(basket_id, product_id)
	prevStock := returnProductStock(product_id)

	if prevStock < product_count {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusMethodNotAllowed,
			"message": "Not enough stock ",
		})
		return
	}

	if prevId != "" { // update product count in Product Basket(PB) if already this product can be found in the basket
		_, err := dbConnect.Model(&entities.BasketProduct{}).
			Set("product_count = ?", product_count+prevCount).
			Where("basket_product_id = ?", prevId).
			Update()

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
		_, insertError := dbConnect.Model(&entities.BasketProduct{
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
	_ = ChangeProductStock(product_id, product_count)

}

func returnProductStock(p_id string) int {
	var product = &entities.Product{}
	err := dbConnect.Model(product).
		Where("product_id = ?", p_id).
		Select()
	if err != nil {
		return 0
	}

	return product.ProductStock
}

func DiscardBasketProduct(c *gin.Context) {
	basketproduct_id := c.Param("basketproduct_id") // takes bp_id parameter
	var basketProduct entities.BasketProduct
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
		_ = UpdateBPCount(basketproduct_id, prevCount-discardedCount)
	}

}

func ChangeProductStock(product_id string, productDemand int) error {

	product := &entities.Product{ProductID: product_id}
	err := dbConnect.Model(product).WherePK().Select()
	if err != nil {
		return err
	} else {

		productStock := product.ProductStock
		// does any mutexLock kind of thing needed here ?
		if productStock < productDemand { //
			return errors.New(" no sufficient product ")
		} else {
			// product update by id
			_, err := dbConnect.Model(&entities.Product{}).
				Set("product_stock = ?", productStock-productDemand).
				Set("updated_at = ?", time.Now()).
				Where("product_id = ?", product_id).
				Update()
			if err != nil {
				return err
			}
			log.Printf("Product Stock Edited Successfully")
			return nil
		}
	}
}

// returns "Basket Product ID" if duplicate
func AddIfDuplicate(b_id string, p_id string) (string, int) {
	basket := &entities.Basket{BasketID: b_id}
	err := dbConnect.Model(basket).
		Relation("BasketProducts").
		Relation("BasketProducts.Product").
		WherePK().
		Select()
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
	basketProduct := &entities.BasketProduct{BasketProductID: bp_id}
	err := dbConnect.Model(basketProduct).WherePK().Select()
	if err != nil {
		log.Printf("Error while getting the Basket Product , Reason: %v\n", err)
	}
	return basketProduct.ProductCount, basketProduct.ProductID
}

func DeleteBP(bp_id string, c *gin.Context) {
	//basketproduct_id := c.Param("basketproduct_id") // takes bp_id parameter
	myBasketProduct := &entities.BasketProduct{BasketProductID: bp_id}

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

func UpdateBPCount(bp_id string, product_count int) error {
	//basketproduct_id := c.Param("basketproduct_id") // takes bp_id parameter
	myBasketProduct := &entities.BasketProduct{BasketProductID: bp_id}

	_, err := dbConnect.Model(myBasketProduct).
		Set("product_count = ?", product_count).
		WherePK().
		Update()

	if err != nil {
		log.Printf("Error, Reason: %v\n", err)
		return err
	}
	return nil
}
