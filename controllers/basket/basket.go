package basket

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	pg "github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	guuid "github.com/google/uuid"
	entities "github.com/mataliksamil/Go_Bootcamp_Final/entities"
)

// Create Basket Table
func CreateBasketTable(db *pg.DB) error {
	opts := &orm.CreateTableOptions{
		FKConstraints: true,
		IfNotExists:   true,
	}

	createError := db.Model(&entities.Basket{}).CreateTable(opts)
	if createError != nil {
		log.Printf("Error while creating basket table, Reason: %v\n", createError)
		return createError
	}
	log.Printf("Basket table created")
	return nil
}

func CreateBasket(c *gin.Context) {
	var basket entities.Basket
	c.BindJSON(&basket)
	user_id := basket.UserID
	basket_id := guuid.New().String() // 	// this will connect user logic
	total_cost := 0.0                 //default value
	total_vat := 0.0
	basket_status := 1 // default value

	_, insertError := dbConnect.Model(&entities.Basket{
		BasketID:     basket_id,
		UserID:       user_id,
		TotalCost:    total_cost,
		TotalVAT:     total_vat,
		BasketStatus: basket_status,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
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
		"message": "Basket created Successfully",
	})
}

func GetSingleBasket(c *gin.Context) {

	basketId := c.Param("basket_id")
	basket := &entities.Basket{BasketID: basketId}

	err := dbConnect.Model(basket).
		Relation("BasketProducts").
		Relation("BasketProducts.Product").
		WherePK().Select()

	//basket.BasketProducts[].Product

	if err != nil {
		log.Printf("Error while getting a single basket, Reason: %v\n", err)
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Basket not found",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Single Basket",
		"data":    basket,
	})
}
func calculateTotalCost(user_id string) (float64, float64, error) {
	totalCost := 0.0
	totalVAT := 0.0
	var myBasket = &entities.Basket{}
	err := dbConnect.Model(myBasket).
		Relation("BasketProducts").
		Relation("BasketProducts.Product").
		Where("user_id=?", user_id).
		Where("basket_status=?", 1).
		Select()
	if err != nil {
		return 0.0, 0.0, err
	}
	myBasketProduct := myBasket.BasketProducts
	for _, myBP := range myBasketProduct {
		totalCost += myBP.Product.Price * float64(myBP.ProductCount)
		totalVAT += myBP.Product.Price * float64(myBP.ProductCount) * float64(myBP.Product.VatRate+100) / 100
	}

	_, err = dbConnect.Model(myBasket).
		Set("total_vat =?", totalVAT).Set("total_cost =?", totalCost).
		Where("user_id=?", user_id).
		Where("basket_status=?", 1).
		Update()
	if err != nil {
		return 0.0, 0.0, err
	}
	return totalCost, totalVAT, nil
}

func updateTotals(totalCost float64, totalVAT float64, user_id string) error {

	var myBasket = &entities.Basket{}
	err := dbConnect.Model(myBasket).
		Set("total_vat =?", totalVAT).Set("total_cost =?", totalCost).
		Where("user_id=?", user_id).
		Where("basket_status=?", 1).
		Select()
	if err != nil {
		return err
	}
	return nil
}

// this func toggles basket status to 0
// by this way  order become completed
func CompleteTheOrder(c *gin.Context) {
	basket_id := c.Param("basket_id")
	basket := &entities.Basket{}
	c.BindJSON(&basket)
	basket_status := 0

	_, err := dbConnect.Model(basket).
		Set("basket_status = ?", basket_status).
		Set("updated_at=?", time.Now()).
		Where("basket_id = ?", basket_id).Update()

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
		"message": "Order Completed ",
		"data":    basket,
	})
}

func DeleteBasket(c *gin.Context) {
	basket_id := c.Param("basket_id")
	basket := &entities.Basket{BasketID: basket_id}
	_, err := dbConnect.Model(basket).WherePK().Delete()
	//err := dbConnect.Delete(product)
	if err != nil {
		log.Printf("Error while deleting a basket, Reason: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Something went wrong",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Basket deleted successfully",
	})

}
