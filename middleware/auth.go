package middleware

import (
	"auth-service/auth"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(context *gin.Context) {

		// Gets the Authorization header from the HTTP request
		tokenString := context.GetHeader("Authorization")

		// Checks if the JWT token was sent.
		if tokenString == "" {
			context.JSON(401, gin.H{"error": "request does not contain an access token"})
			context.Abort()
			return
		}

		// Validates the JWT token
		err := auth.ValidateToken(tokenString)

		if err != nil {
			context.JSON(401, gin.H{"error": err.Error()})
			context.Abort()
			return
		}

		// Allows middleware to continue with the flow and the request
		context.Next()
	}
}
