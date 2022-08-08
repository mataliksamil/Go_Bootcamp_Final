package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mataliksamil/Go_Bootcamp_Final/controllers"
)

func Routes(router *gin.Engine) {
	router.GET("/", welcome)
	router.GET("/user", controllers.GetAllUsers)
	router.POST("/user", controllers.CreateUser)
	//router.GET("/user/:userId", controllers.GetSingleUser)
	router.PUT("/user/:userId", controllers.EditUserName)
	router.DELETE("/user/:userId", controllers.DeleteUser)
	router.GET("/user/:userId", controllers.GetUsersAllBaskets)

	router.GET("/user/1/:userId", controllers.GetUsersActiveBasket)

	router.GET("/product", controllers.GetAllProducts)
	router.POST("/product", controllers.CreateProduct)
	router.GET("/product/:product_id", controllers.GetSingleProduct)
	router.PUT("/product/:product_id", controllers.EditProductStock)
	router.DELETE("/product/:product_id", controllers.DeleteProduct)
	//router.POST("/product/:product_id/:basket_id", controllers.ProductToBasket)

	router.GET("/basket", controllers.GetAllProducts)
	router.POST("/basket", controllers.CreateBasket)
	router.GET("/basket/:basket_id", controllers.GetSingleBasket)
	//router.GET("/basket/:basket_owner_id",controllers.CreateBasket)
	router.PUT("/basket/:basket_id", controllers.CompleteTheOrder)
	router.DELETE("/basket/:basket_id", controllers.DeleteBasket)

	router.POST("/basket_product", controllers.AddProductToBasket)
	router.PUT("/basket_product/:basketproduct_id", controllers.DiscardBasketProduct)
}
func welcome(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "Welcome To API",
	})
}
