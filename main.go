package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/harimadhav/event-management/controllers"
	"github.com/harimadhav/event-management/routes"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	godotenv.Load()

	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	db := client.Database("eventdb")
	userCollection := db.Collection("users")
	controllers.InitAuth(userCollection)

	r := gin.Default()
	routes.AuthRoutes(r)
	r.Run(":" + os.Getenv("PORT"))
}
