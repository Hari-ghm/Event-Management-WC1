package controllers

import (
    "context"
    "net/http"
    "time"
    "fmt"
    "github.com/gin-gonic/gin"
    "github.com/Hari-ghm/Event-Management-WC1/models"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
)

var eventCollection *mongo.Collection

func InitEvent(collection *mongo.Collection) {
    eventCollection = collection
}

func CreateEvent(c *gin.Context) {
    var event models.Event
    if err := c.BindJSON(&event); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    email, _ := c.Get("email")
    event.ID = primitive.NewObjectID()
    event.CreatedBy = email.(string)

    email, ok := c.Get("email")
    if !ok {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "email not found in context"})
        return
    }
    fmt.Println("Creating event for:", email)

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    _, err := eventCollection.InsertOne(ctx, event)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create event"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Event created"})

    
}

func ListEvents(c *gin.Context) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    cursor, err := eventCollection.Find(ctx, bson.M{})
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch events"})
        return
    }
    defer cursor.Close(ctx)

    var events []models.Event
    if err := cursor.All(ctx, &events); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading events"})
        return
    }

    c.JSON(http.StatusOK, events)
}

func GetEventByID(c *gin.Context) {
    id := c.Param("id")
    objID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
        return
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    var event models.Event
    err = eventCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&event)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
        return
    }

    c.JSON(http.StatusOK, event)
}

func UpdateEvent(c *gin.Context) {
    id := c.Param("id")
    objID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
        return
    }

    var updateData models.Event
    if err := c.BindJSON(&updateData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    email, _ := c.Get("email") // Authenticated user

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // Only allow owner to update
    filter := bson.M{"_id": objID, "created_by": email}
    update := bson.M{"$set": bson.M{
        "title":       updateData.Title,
        "description": updateData.Description,
        "date":        updateData.Date,
        "location":    updateData.Location,
    }}

    result, err := eventCollection.UpdateOne(ctx, filter, update)
    if err != nil || result.ModifiedCount == 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Update failed or not authorized"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Event updated"})
}

func DeleteEvent(c *gin.Context) {
    id := c.Param("id")
    objID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
        return
    }

    email, _ := c.Get("email")

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    result, err := eventCollection.DeleteOne(ctx, bson.M{"_id": objID, "created_by": email})
    if err != nil || result.DeletedCount == 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Delete failed or not authorized"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Event deleted"})
}
