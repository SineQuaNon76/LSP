package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"GoProject/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// jwtKey is the secret key used for JWT signing.
var jwtKey = []byte("token")

// GenerateToken generates a new JWT token for a user with a specific username.
func GenerateToken(user models.User) (string, error) {
	// Token is valid for 24 hours.

	claims := jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}

	// Create a new JWT token with the specified claims.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key and get the complete signed token string.
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// VerifyToken verifies the given JWT token and returns the claims if valid.
func VerifyToken(tokenString string) (*models.User, error) {
	// Parse the token with the Claims structure.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	// Check if the token is valid and extract the claims.
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		id := int(claims["id"].(float64)) // Adjust for the type of "id" in your claims
		username := claims["username"].(string)

		user := models.User{
			ID:       id,
			Username: username,
		}
		return &user, nil
	}

	return nil, errors.New("invalid token")
}

// JWTMiddleware is middleware to secure routes with JWT.
func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the token from the cookie
		tokenString, err := c.Cookie("token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak ada JWT"})
			c.Abort()
			return
		}

		fmt.Println("Token received:", tokenString)

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak ada JWT"})
			c.Abort()
			return
		}

		// Verify the token and get the claims.
		claims, err := VerifyToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid"})
			c.Abort()
			return
		}

		// Set the "username" in the context for further processing.
		c.Set("username", claims.Username)

		c.Next()
	}
}
