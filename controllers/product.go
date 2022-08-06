package controllers

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	pg "github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

type Product struct {
	ProductID    string  `pg:",pk" json:"product_id"`
	ProductName  string  `pg:"uniqe" json:"product_name"`
	ProductStock int     `json:"product_stock"`
	Price        float64 `json:"price"`
	VatRate      int     `json:"vat_rate"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

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
	/* insertError := dbConnect.Insert(&Product{
		Product_ID:    product_id,
		Product_Name:  product_name,
		Product_Stock: product_stock,
		Price:         price,
		Vat_Rate:      vat_rate,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}) */
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
	_, err := dbConnect.Model(&Product{}).Set("product_stock = ?", product_stock).Where("id = ?", product_id).Update()
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

/* func createNewPartition(db *pg.DB, currentTime time.Time) error {
	firstOfMonth := time.Date(currentTime.Year(), currentTime.Month(), 1, 0, 0, 0, 0, time.UTC)
	firstOfNextMonth := firstOfMonth.AddDate(0, 1, 0)

	year := firstOfMonth.Format("2006")
	month := firstOfMonth.Format("01")
	sql := fmt.Sprintf(
		`CREATE TABLE IF NOT EXISTS logs_y%s_m%s PARTITION OF logs FOR VALUES FROM ('%s') TO ('%s');`,
		year, month,
		firstOfMonth.Format(time.RFC3339Nano),
		firstOfNextMonth.Format(time.RFC3339Nano),
	)

	_, err := db.Exec(sql)
	return err
}
*/

func ProductToBasket(c *gin.Context) {
	/* 	product_id := c.Param("product_id")
	   	basket_id := c.Param("basket_id")

	   	product := new(Product)
	   	c.BindJSON(&product)
	   	quantity := product.Product_Stock
	   	log.Printf("%#v", product)
	   	basket_product := &BasketProduct{
	   		Basket_ID:     basket_id,
	   		Product_ID:    product_id,
	   		Product_Count: quantity,
	   	}

	   	_, err := dbConnect.Model(basket_product).Insert()
	   	if err != nil {
	   		panic(err)
	   	} */
	/*
		_, err = dbConnect.Model(&Product{}).Set("product_stock = ?", product_stock).Where("id = ?", product_id).Update()
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
	*/
}

/*
func getProductStock(product_id string) {

	product_id := c.Param("product_id")
	product := &Product{Product_ID: product_id}
	err := dbConnect.Model(product).Query()
	err := dbConnect.Model(product).WherePK()
	if err != nil {
		log.Printf("Error while getting a single product, Reason: %v\n", err)
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
*/
