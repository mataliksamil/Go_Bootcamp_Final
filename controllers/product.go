package controllers

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	pg "github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

// Create Product Table
func CreateProductTable(db *pg.DB) error {
	opts := &orm.CreateTableOptions{
		FKConstraints: true,
		IfNotExists:   true,
	}

	createError := db.Model(&Product{}).CreateTable(opts)
	if createError != nil {
		log.Printf("Error while creating Product table, Reason: %v\n", createError)
		return createError
	}
	log.Printf("Product table created")
	return nil
}

func GetAllProducts(c *gin.Context) {
	var products []Product

	err := dbConnect.Model(&products).Select()
	if err != nil {
		log.Printf("Error while getting all products, Reason: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Something went wrong",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "All Users",
		"data":    products,
	})

}

func CreateProduct(c *gin.Context) {
	var product Product
	c.BindJSON(&product)
	product_name := product.ProductName
	product_id := product.ProductID
	product_stock := product.ProductStock
	price := product.Price
	vat_rate := product.VatRate
	_, insertError := dbConnect.Model(&Product{
		ProductID:    product_id,
		ProductName:  product_name,
		ProductStock: product_stock,
		Price:        price,
		VatRate:      vat_rate,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}).Insert()

	if insertError != nil {
		log.Printf("Error while inserting new Product into db, Reason: %v\n", insertError)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Something went wrong",
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"message": "Product created Successfully",
	})

}

func GetSingleProduct(c *gin.Context) {
	product_id := c.Param("product_id")
	product := &Product{ProductID: product_id}

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

}

func EditProductStock(c *gin.Context) {
	product_id := c.Param("product_id")
	var product Product
	c.BindJSON(&product)
	product_stock := product.ProductStock

	_, err := dbConnect.Model(&Product{}).
		Set("product_stock = ?", product_stock).
		Where("product_id = ?", product_id).
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
		"message": "Product Stock Edited Successfully",
	})

}

func DeleteProduct(c *gin.Context) {
	product_id := c.Param("product_id")
	product := &Product{ProductID: product_id}
	_, err := dbConnect.Model(product).WherePK().Delete()
	//err := dbConnect.Delete(product)
	if err != nil {
		log.Printf("Error while deleting a product, Reason: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Something went wrong",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Product deleted successfully",
	})

}

func ChangeProductStock(product_id string, productDemand int) error {

	product := &Product{ProductID: product_id}
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
			_, err := dbConnect.Model(&Product{}).
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
