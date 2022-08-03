package main

import (
	"log"

	"github.com/gin-gonic/gin"
	routes "github.com/mataliksamil/Go_Bootcamp_Final/routes"
)

func main() {
	// Init Router
	router := gin.Default()
	// Route Handlers / Endpoints
	routes.Routes(router)
	log.Fatal(router.Run(":8080"))
}
