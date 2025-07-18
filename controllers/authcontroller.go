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

func RegisterUser(c *gin.Context) {
	var user models.User

	// binding json to struct
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// first create a if
	user.ID = primitive.NewObjectID()

	// getting collection from the database
	userCollection := config.Client.Database("Gin").Collection("user")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := userCollection.InsertOne(ctx, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"userId":  result.InsertedID,
	})
}
