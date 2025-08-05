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

func CreatedProduct(c *gin.Context) {
   var product models.Product
   err := c.ShouldBindJSON(&product)
   if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
   }

   product.ID = primitive.NewObjectID()
   productCollection := config.Client.Database("Gin").Collection("product")
   ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
   defer cancel()
   
   result, err := productCollection.InsertOne(ctx, product)
   if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create the product"})
		return
   }

   c.JSON(http.StatusCreated, gin.H{
	"message":"Product is created successfully",
	"data":result,
   })
}