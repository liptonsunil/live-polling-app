package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Global database clients
var (
	mongoClient *mongo.Client
	redisClient *redis.Client
	ctx         = context.Background()
)

func initDatabases() {
	// Connect to MongoDB
	mongoOptions := options.Client().ApplyURI("mongodb+srv://leofrancis895_db_user:<L3GV1oGb74N2mgnH>@cluster10.knijtlk.mongodb.net/?appName=Cluster10")
	client, err := mongo.Connect(ctx, mongoOptions)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	mongoClient = client
	fmt.Println("Connected to MongoDB!")

// Connect to Upstash Redis
	opt, err := redis.ParseURL("rediss://default:gQAAAAAABGYbAAIgcDFjNzY3ZTIyZDY3NDI0NGVkOTdmYTk5YTAyYWFhMTY5Ng@first-dove-288283.upstash.io:6379")
	if err != nil {
		log.Fatal("Failed to parse Redis URL:", err)
	}
	redisClient = redis.NewClient(opt)

	// Test Redis connection
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}
	fmt.Println("Connected to Redis!")
	// Test Redis connection
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}
	fmt.Println("Connected to Redis!")
}

func main() {
	initDatabases()

	// Initialize the Gin router
	r := gin.Default()

	// Health check route
	r.GET("/ping", func(c *gin.Context) {c.JSON(200, gin.H{"message": "pong"})
	})

        r.Use(cors.New(cors.Config{
                AllowOrigins:     []string{"*"}, // Allows all origins
                AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
                AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
	}))

	// Poll CRUD routes
	// Authentication Routes
	r.POST("/api/signup", Signup)
	r.POST("/api/login", Login)

	// Protected Poll Route
	r.POST("/api/polls", AuthMiddleware(), CreatePoll)
	r.GET("/api/polls/:id", GetPoll)
        // Real-time voting routes
	r.POST("/api/polls/:id/vote", SubmitVote)
	r.GET("/ws/polls/:id", WatchPoll)
	// Run the server on port 8080
	fmt.Println("Server is running on http://localhost:8080")
	r.Run(":8080")
}
