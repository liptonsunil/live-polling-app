package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreatePoll handles POST /api/polls
func CreatePoll(c *gin.Context) {
	var poll Poll
	
	// Validate input on the backend[cite: 1]
	if err := c.ShouldBindJSON(&poll); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	collection := mongoClient.Database("live_polling").Collection("polls")
	result, err := collection.InsertOne(context.Background(), poll)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save poll"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Poll created",
		"id":      result.InsertedID,
	})
}

// GetPoll handles GET /api/polls/:id
func GetPoll(c *gin.Context) {
	pollID := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(pollID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid poll ID"})
		return
	}

	var poll Poll
	collection := mongoClient.Database("live_polling").Collection("polls")
	
	err = collection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&poll)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Poll not found"})
		return
	}

	c.JSON(http.StatusOK, poll)
}
