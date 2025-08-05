package main

import (
	"example/go-gin-resume-sender/config"
	"example/go-gin-resume-sender/routes"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

// func for getting all the albums
func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}
func main() {
	db := config.ConnectToDB()
	collection := db.Database("Gin").Collection("user")
	fmt.Print("the collection", collection)
	router := gin.Default()

	router.GET(`/`, func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	router.GET("/albums", getAlbums)
	routes.AuthRoutes(router)
	routes.CategoryRoutes(router)
	routes.ProductRoutes(router)
	router.Run("localhost:8000")
}
