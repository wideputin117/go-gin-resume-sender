package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func ConnectToDB() *mongo.Client {
	// Load environment variables
	_ = godotenv.Load()
	// if err != nil {
	// 	log.Fatal("❌ Unable to load environment variables")
	// }

	// Get the MongoDB URI from .env
	mongoURL := os.Getenv("MONGODB_URI")
	if mongoURL == "" {
		log.Fatal("❌ MONGODB_URI not found in environment")
	}

	// Create a new client
	clientOptions := options.Client().ApplyURI(mongoURL)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("❌ MongoDB connection failed:", err)
	}

	// Ping the database
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("❌ MongoDB ping failed:", err)
	}

	fmt.Println("✅ Connected to MongoDB successfully!")
	Client = client
	return Client
}
