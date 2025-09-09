package main

import (
	"example/go-gin-resume-sender/config"
	"example/go-gin-resume-sender/routes"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	fmt.Print(port)
	if port == "" {
    port = "5000" // default fallback for local dev
}
	db := config.ConnectToDB()
	collection := db.Database("Gin").Collection("user")
	fmt.Print("the collection", collection)
	router := gin.Default()

	router.GET(`/`, func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	routes.AuthRoutes(router)
	routes.CategoryRoutes(router)
	routes.ProductRoutes(router)
	routes.CompanyRoutes(router)
	routes.ResumeRoutes(router)
	router.Run(":" + port)
}
