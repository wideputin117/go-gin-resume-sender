package controllers

import (
	"context"
	"example/go-gin-resume-sender/config"
	"example/go-gin-resume-sender/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreatedCategory(c *gin.Context) {
	var category models.Category
    
	err := c.ShouldBindJSON(&category)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	category.ID = primitive.NewObjectID()

	categoryCollection := config.Client.Database("Gin").Collection("category")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := categoryCollection.InsertOne(ctx, category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":"Failed to create te category",
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message":"Category added",
		"data":result,
	})
}