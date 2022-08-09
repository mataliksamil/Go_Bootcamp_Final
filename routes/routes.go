package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	basket "github.com/mataliksamil/Go_Bootcamp_Final/controllers/basket"
	basketProduct "github.com/mataliksamil/Go_Bootcamp_Final/controllers/basketProduct"
	product "github.com/mataliksamil/Go_Bootcamp_Final/controllers/product"
	user "github.com/mataliksamil/Go_Bootcamp_Final/controllers/user"
)

func Routes(router *gin.Engine) {
	router.GET("/", welcome)
	router.GET("/user", user.GetAllUsers)
	router.POST("/user", user.CreateUser)
	//router.GET("/user/:userId", controllers.GetSingleUser)
	router.PUT("/user/:userId", user.EditUserName)
	router.DELETE("/user/:userId", user.DeleteUser)

	router.GET("/user/1/:userId", user.GetUsersAllBaskets)
	router.GET("/user/:userId", user.GetUsersActiveBasket)

	router.GET("/product", product.GetAllProducts)
	router.POST("/product", product.CreateProduct)
	router.GET("/product/:product_id", product.GetSingleProduct)
	router.PUT("/product/:product_id", product.EditProductStock)
	router.DELETE("/product/:product_id", product.DeleteProduct)
	//router.POST("/product/:product_id/:basket_id", controllers.ProductToBasket)

	router.GET("/basket", product.GetAllProducts)
	router.POST("/basket", basket.CreateBasket)
	router.GET("/basket/:basket_id", basket.GetSingleBasket)
	//router.GET("/basket/:basket_owner_id",controllers.CreateBasket)
	router.PUT("/basket/:basket_id", basket.CompleteTheOrder)
	router.DELETE("/basket/:basket_id", basket.DeleteBasket)

	router.POST("/basket_product", basketProduct.AddProductToBasket)
	router.PUT("/basket_product/:basketproduct_id", basketProduct.DiscardBasketProduct)

}
func welcome(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "Welcome To API",
	})
}
