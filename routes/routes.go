package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	controllers "github.com/mataliksamil/Go_Bootcamp_Final/controllers"
)

func Routes(router *gin.Engine) {
	router.GET("/", welcome)
	router.GET("/user", controllers.GetAllUsers)
	router.POST("/user", controllers.CreateUser)
	router.GET("/user/:userId", controllers.GetSingleUser)
	router.PUT("/user/:userId", controllers.EditUser)
	router.DELETE("/user/:userId", controllers.DeleteUser)
}
func welcome(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "Welcome To API",
	})
}
