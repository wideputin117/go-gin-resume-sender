package routes

import (
	"example/go-gin-resume-sender/controllers"

	"github.com/gin-gonic/gin"
)

func CompanyRoutes(c *gin.Engine) {
	companycontroller := controllers.NewCompanyController()
	company := c.Group("/api/v1/company")
	{
		company.POST("/", companycontroller.CreateCompany)
		company.GET("/", companycontroller.GetCompanies)

		companyId := company.Group("/:companyId")
		{
			companyId.GET("", companycontroller.GetSingleCompany)
			companyId.DELETE("", companycontroller.DeleteCompany)
		}
	}
}
