package controllers

import (
	"context"
	"example/go-gin-resume-sender/config"
	"example/go-gin-resume-sender/models"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

func GetProducts(c *gin.Context) {
	productCollection := config.Client.Database("Gin").Collection("product")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// preparing the pipeline for performing the aggregation
 	pipeline := mongo.Pipeline{
		{
			{"$lookup", bson.D{
				{Key: "from", Value: "category"},
				{Key: "localField", Value: "category"},
				{Key: "foreignField", Value: "_id"},
				{Key: "as", Value: "categoryInfo"},
			}},
		},
		{
			{"$unwind", bson.D{
				{"path", "$categoryInfo"},
				{"preserveNullAndEmptyArrays", true},  
			}},
		},
	}

	cursor, err := productCollection.Aggregate(ctx, pipeline)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}
	defer cursor.Close(ctx)

	var products []bson.M
	if err = cursor.All(ctx, &products); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode products"})
		return
	}

	if len(products) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No products found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Products fetched successfully",
		"products": products,
	})
}

func GetSingleProduct(c *gin.Context){
	id := c.Param("productId")
	fmt.Println("the id is", id)
    var product models.Product
	if id == ""{
		c.JSON(http.StatusBadRequest, gin.H{"error":"Inavlid Product Id"})
		return
	}
	objId, err:= primitive.ObjectIDFromHex(id)
	productCollection :=config.Client.Database("Gin").Collection("product")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = productCollection.FindOne(ctx,bson.M{"_id":objId}).Decode(&product)

	if err != nil{
		c.JSON(http.StatusNotFound,gin.H{"message":"No product is found"})
	    return
	}
	c.JSON(http.StatusFound, gin.H{"message": "Product is Found", "data": product})
} 

func DeleteSingleProduct(c *gin.Context){
	id := c.Param("productId")
    var product models.Product

	productCollection := config.Client.Database("Gin").Collection("product")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

    objectId ,err := primitive.ObjectIDFromHex(id)
    if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error":""})
        return
	}
	 err = productCollection.FindOneAndDelete(ctx,bson.M{"_id":objectId}).Decode(&product)
	 if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error":"Unable to delete the product"})
		return 
	 }
	 c.JSON(http.StatusCreated, gin.H{"message":"Successfully deleted the product","success":true, "data":product})
}