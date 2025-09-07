package routes

import (
	"example/go-gin-resume-sender/controllers"

	"github.com/gin-gonic/gin"
)

func ResumeRoutes(c *gin.Engine) {
	resume := c.Group("/api/v1/resume")
	{
		resume.POST("/send", controllers.ParseResume)
	}

}
