package controllers

import (
	"context"
	"example/go-gin-resume-sender/config"
	"example/go-gin-resume-sender/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CompanyController struct {
	Collection *mongo.Collection
}

// CONTRUCTOR FUNCTION
func NewCompanyController() *CompanyController{
	return &CompanyController{
		Collection: config.Client.Database("Gin").Collection("company"),
	}
}


func (company *CompanyController) CreateCompany(c *gin.Context){
    var newCompany models.Company
    err := c.ShouldBindJSON(&newCompany)
	
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error":"Invalid data"})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(),10*time.Second)
	defer cancel()
    newCompany.ID= primitive.NewObjectID()
	result , err := company.Collection.InsertOne(ctx, newCompany)
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error":"Failed to create the company"})
		return
	}

	c.JSON(http.StatusCreated,gin.H{"message":"Company Created","data":result})
}