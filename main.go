package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/Hari-ghm/Event-Management-WC1/controllers"
	"github.com/Hari-ghm/Event-Management-WC1/routes"
	"github.com/Hari-ghm/Event-Management-WC1/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Load .env BEFORE accessing any env vars
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize JWT secret
	utils.InitJWTSecret()

	// Initialize Gin router
	r := gin.Default()

	// Connect to MongoDB
	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		log.Fatal("Mongo client creation failed:", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal("Mongo connection failed:", err)
	}

	// Get DB and user collection
	db := client.Database("eventdb")
	userCollection := db.Collection("users")
	controllers.InitAuth(userCollection)

	// Set up routes (ensure r is passed, not created inside)
	routes.SetupRoutes(r)

	// Start server on port (e.g. 8080)
	err = r.Run(":" + os.Getenv("PORT"))
	if err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
