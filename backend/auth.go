package main

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("super_secret_internship_key")

// Signup hashes the password and saves the user to MongoDB
func Signup(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	user.Password = string(hashed)

	collection := mongoClient.Database("live_polling").Collection("users")
	_, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}

// Login verifies the password and returns a JWT
func Login(c *gin.Context) {
	var input User
	c.ShouldBindJSON(&input)

	var dbUser User
	collection := mongoClient.Database("live_polling").Collection("users")
	err := collection.FindOne(context.Background(), bson.M{"username": input.Username}).Decode(&dbUser)

	if err != nil || bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(input.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": dbUser.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, _ := token.SignedString(jwtSecret)

	c.JSON(http.StatusOK, gin.H{"token": tokenString, "username": dbUser.Username})
}

// AuthMiddleware intercepts requests and blocks them if the token is missing/invalid
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authorization token"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, _ := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("username", claims["username"])
			c.Next() // Token is valid, allow the request to proceed
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
		}
	}
}
