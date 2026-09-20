package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true 
	},
}

func SubmitVote(c *gin.Context) {
	pollID := c.Param("id")
	
	var req struct {
		OptionID string `json:"option_id"`
		VoterID  string `json:"voter_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	objID, err := primitive.ObjectIDFromHex(pollID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid poll ID"})
		return
	}

	collection := mongoClient.Database("live_polling").Collection("polls")
	
	filter := bson.M{
		"_id": objID, 
		"options.id": req.OptionID,
		"voted_by": bson.M{"$ne": req.VoterID}, 
	}
	
	update := bson.M{
		"$inc": bson.M{"options.$.votes": 1},
		"$push": bson.M{"voted_by": req.VoterID},
	}
	
	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to vote"})
		return
	}
	
	if result.ModifiedCount == 0 {
		c.JSON(http.StatusForbidden, gin.H{"error": "You have already voted!"})
		return
	}

	var updatedPoll Poll
	collection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&updatedPoll)

	pollJSON, _ := json.Marshal(updatedPoll)
	channelName := "poll_updates_" + pollID
	redisClient.Publish(context.Background(), channelName, pollJSON)

	c.JSON(http.StatusOK, gin.H{"message": "Vote registered"})
}

func WatchPoll(c *gin.Context) {
	pollID := c.Param("id")
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer ws.Close()

	channelName := "poll_updates_" + pollID
	pubsub := redisClient.Subscribe(context.Background(), channelName)
	defer pubsub.Close()

	ch := pubsub.Channel()

	for msg := range ch {
		err := ws.WriteMessage(websocket.TextMessage, []byte(msg.Payload))
		if err != nil {
			break 
		}
	}
}
