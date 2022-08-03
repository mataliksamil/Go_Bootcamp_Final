package main

import (
	"log"

	"github.com/gin-gonic/gin"

	config "github.com/mataliksamil/Go_Bootcamp_Final/config"
	routes "github.com/mataliksamil/Go_Bootcamp_Final/routes"
)

func main() {
	config.Connect()
	// Init Router
	router := gin.Default()
	// Route Handlers / Endpoints
	routes.Routes(router)
	log.Fatal(router.Run(":8081"))
}
