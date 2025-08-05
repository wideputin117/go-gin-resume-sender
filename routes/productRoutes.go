package routes

import (
	"example/go-gin-resume-sender/controllers"

	"github.com/gin-gonic/gin"
)

func ProductRoutes(c *gin.Engine){
	product := c.Group("/api/v1/product")
	{
		product.POST("/create", controllers.CreatedProduct)
	}
}