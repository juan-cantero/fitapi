package middleware

import (
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthRequired is a middleware that validates JWT tokens from Supabase
// It extracts the token from the Authorization header and validates it
// If valid, it stores the user_id in the Gin context for handlers to use
func AuthRequired() gin.HandlerFunc {
	// Get JWT secret from environment
	jwtSecret := os.Getenv("SUPABASE_JWT_SECRET")
	if jwtSecret == "" {
		panic("SUPABASE_JWT_SECRET environment variable is required")
	}

	return func(c *gin.Context) {
		// 1. Extract Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{
				"error": "missing authorization header",
			})
			c.Abort()
			return
		}

		// 2. Extract token (remove "Bearer " prefix)
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			// "Bearer " prefix not found
			c.JSON(401, gin.H{
				"error": "invalid authorization header format, expected 'Bearer <token>'",
			})
			c.Abort()
			return
		}

		// 3. Parse and validate JWT token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(401, gin.H{
				"error": "invalid or expired token",
			})
			c.Abort()
			return
		}

		// 4. Extract claims (user information)
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(401, gin.H{
				"error": "invalid token claims",
			})
			c.Abort()
			return
		}

		// 5. Extract user_id (sub claim) and email
		userID, ok := claims["sub"].(string)
		if !ok {
			c.JSON(401, gin.H{
				"error": "invalid user_id in token",
			})
			c.Abort()
			return
		}

		email, _ := claims["email"].(string) // Optional

		// 6. Store user information in context for handlers to use
		c.Set("user_id", userID)
		c.Set("user_email", email)

		// 7. Continue to the next handler
		c.Next()
	}
}
