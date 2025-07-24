package controllers

import (
	"context"
	"example/go-gin-resume-sender/config"
	"example/go-gin-resume-sender/models"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)
 type CookieOptions struct {
	MaxAge   int
	Path     string
	Domain   string
	Secure   bool
	HttpOnly bool
	SameSite http.SameSite // Use the standard library's SameSite type
}


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

func LoginUser (c *gin.Context){
	err := godotenv.Load()
	if err != nil {
		log.Fatal("❌ Unable to load environment variables")
	}
   
	 
	type Login struct{
		Email  string
		Password string
	}
	var login Login
  err = c.ShouldBindJSON(&login)

  if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
  }
  
  usercollection := config.Client.Database("Gin").Collection("user")
  var user models.User
  ctx , cancel := context.WithTimeout(context.Background(),5*time.Second)
  defer cancel()
  err = usercollection.FindOne(ctx, bson.M{"email":login.Email}).Decode(&user)
   if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching user"})
		}
		return
	}
   
	if login.Password != user.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error":"Password do not match"})
	    return
	}
	var (
		key []byte
		t   *jwt.Token
		s   string
      )
	jwt_token := os.Getenv("JWT_SECRET")
    if jwt_token == ""{
		log.Fatal("Secret not found")
	}
    
	t= jwt.NewWithClaims(jwt.SigningMethodHS256,jwt.MapClaims{
		  "iss": "my-auth-server", 
			"sub": "john", 
			"foo": 2,
	})
    s,err =t.SignedString(key)
    if err != nil {
        log.Fatal("Error in signing the key")
	}
	c.Header("Authorization", "Bearer "+s)

    c.JSON(http.StatusOK, gin.H{"message": "Login Successful"})

}