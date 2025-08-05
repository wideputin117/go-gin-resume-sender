package routes

import (
	"example/go-gin-resume-sender/controllers"

	"github.com/gin-gonic/gin"
)

func CategoryRoutes(c *gin.Engine) {
   category := c.Group("/api/v1/category")
   {
	category.POST("/create", controllers.CreatedCategory)
   }
}