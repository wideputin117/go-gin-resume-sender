package routes

import (
	"example/go-gin-resume-sender/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(c *gin.Engine) {
	auth := c.Group("/api/v1/auth")
	{
		auth.POST("/signup", controllers.RegisterUser)
		auth.POST("/login",controllers.LoginUser)

 	}
 }
