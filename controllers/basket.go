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

type Basket struct {
	BasketID     string  `pg:",pk" json:"basket_id"`
	BasketName   string  `json:"basket_name"`
	TotalCost    float64 `json:"total_cost"`
	BasketStatus int     `json:"basket_status"`

	BasketProducts []*BasketProduct `pg:"rel:has-many"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Create Basket Table
func CreateBasketTable(db *pg.DB) error {
	opts := &orm.CreateTableOptions{
		FKConstraints: true,
		IfNotExists:   true,
	}

	createError := db.Model(&Basket{}).CreateTable(opts)
	if createError != nil {
		log.Printf("Error while creating basket table, Reason: %v\n", createError)
		return createError
	}
	log.Printf("Basket table created")
	return nil
}

func CreateBasket(c *gin.Context) {
	var basket Basket
	c.BindJSON(&basket)

	basket_id := guuid.New().String() // 	// this will connect user logic
	total_cost := 0.0                 //default value
	basket_status := 1                // default value
	basket_name := basket.BasketName
	_, insertError := dbConnect.Model(&Basket{
		BasketID:     basket_id,
		BasketName:   basket_name,
		TotalCost:    total_cost,
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
	/* basket_id := c.Param("basket_id")
	basket := &Basket{Basket_ID: basket_id} */

	basketId := c.Param("basket_id")
	basket := &Basket{BasketID: basketId}

	err := dbConnect.Model(basket).Relation("BasketProducts").Relation("BasketProducts.Product").WherePK().Select()
	if err != nil {
		log.Printf("%#v", basket)
	}
	/* 	if err != ni
	   		log.Printf("Error while getting a single basket, Reason: %v\n", err)
	   		c.JSON(http.StatusNotFound, gin.H{
	   			"status":  http.StatusNotFound,
	   			"message": "Basket not found",
	   		})
	   		return
	   	}
	   	c.JSON(http.StatusOK, gin.H{
	   		"status":  http.StatusOK,
	   		"message": "Single Baskeet",
	   		"data":    basket,
	   	}) */

}

func EditBasketStatus(c *gin.Context) {
	basket_id := c.Param("product_id")
	var basket Basket
	c.BindJSON(&basket)
	basket_status := basket.BasketStatus
	_, err := dbConnect.Model(&Basket{}).Set("basket_status = ?", basket_status).Where("basket_id = ?", basket_id).Update()
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
		"message": "Basket Status Edited Successfully",
	})

}

func DeleteBasket(c *gin.Context) {
	basket_id := c.Param("basket_id")
	basket := &Basket{BasketID: basket_id}
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
