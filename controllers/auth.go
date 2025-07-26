package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/harimadhav/event-management-wc1/models"
	"github.com/harimadhav/event-management/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection

func InitAuth(collection *mongo.Collection) {
	userCollection = collection
}

func Register(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	count, _ := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	user.Password = string(hash)

	user.ID = primitive.NewObjectID()
	_, err := userCollection.InsertOne(ctx, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Registered successfully"})
}

func Login(c *gin.Context) {
	var loginData models.User
	if err := c.BindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := userCollection.FindOne(ctx, bson.M{"email": loginData.Email}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Wrong password"})
		return
	}

	token, err := utils.GenerateToken(user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
