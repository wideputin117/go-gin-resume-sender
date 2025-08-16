package controllers

import (
	"context"
	"example/go-gin-resume-sender/config"
	"example/go-gin-resume-sender/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CompanyController struct {
	Collection *mongo.Collection
}

// CONTRUCTOR FUNCTION
func NewCompanyController() *CompanyController {
	return &CompanyController{
		Collection: config.Client.Database("Gin").Collection("company"),
	}
}

func (company *CompanyController) CreateCompany(c *gin.Context) {
	var newCompany models.Company
	err := c.ShouldBindJSON(&newCompany)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	newCompany.ID = primitive.NewObjectID()
	result, err := company.Collection.InsertOne(ctx, newCompany)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create the company"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Company Created", "data": result})
}

func (company *CompanyController) GetCompanies(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, err := company.Collection.Find(ctx, primitive.M{})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error in getting the companies"})
		return
	}
	var companies []models.Company
	if err = result.All(ctx, &companies); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding companies"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": companies})
}

func (company *CompanyController) DeleteCompany(c *gin.Context) {
	var companyId = c.Param("companyId")
	var company_data models.Company
	objectId, err := primitive.ObjectIDFromHex(companyId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to convert the bson id"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = company.Collection.FindOneAndDelete(ctx, bson.M{"_id": objectId}).Decode(&company_data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "No Product Found"})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"data": company_data})
}

func (company *CompanyController) GetSingleCompany(c *gin.Context) {
	id := c.Param("companyId")
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "Bad objecy id"})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var company_data models.Company
	err = company.Collection.FindOne(ctx, bson.M{"_id": objectId}).Decode(&company_data)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete"})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"message": "Successfullt Returned", "success": true, "data": company_data})
}
