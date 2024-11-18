package controllers

import (
	"auth-service/auth"
	"auth-service/database"
	"auth-service/models"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TokenRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func GenerateToken(context *gin.Context) {
	var request TokenRequest
	var user models.User

	if err := context.ShouldBindBodyWithJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	// Checks if the email and password exist and fetches the first matching record
	record := database.Instance.Where("email = ?", request.Email).First(&user)

	// Checks for possible errors on query
	if record.Error != nil {
		// If a  record is not found, logs the error
		if !errors.Is(record.Error, gorm.ErrRecordNotFound) {
			log.Printf("Database error: %v", record.Error)
		}

		// Return a 401 code to prevent enumeration attacks
		context.JSON(http.StatusUnauthorized, gin.H{"error": "invalide username or password"})
		context.Abort()
		return
	}

	// Validates password
	credentialError := user.CheckPassword(request.Password)
	if credentialError != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "invalide username or password"})
		context.Abort()
		return
	}

	// Generates JWT token and checks for possible errors
	accessToken, err := auth.GenerateJWT(user.Email, user.Username)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	// Generates JWT refresh token and checks for possible errors
	refreshToken, err := auth.GenerateRefreshToken(user.Email, user.Username)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	// Stores refresh token in a cookie for increased security
	context.SetCookie("refresh_token", refreshToken, 7*24*3600, "/", "localhost", true, true)
	context.JSON(http.StatusOK, gin.H{"access_token": accessToken})
}

func GenerateRefreshToken(context *gin.Context) {
	// Fetches the refresh token from the "refresh_token" cookie
	refreshToken, err := context.Cookie("refresh_token")

	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token required"})
		context.Abort()
		return
	}

	// Validates refresh token and retrieves claims
	claims, err := auth.ValidateToken(refreshToken)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	// Generates a new access_code with the previously fetched claims
	newAccessToken, err := auth.GenerateJWT(claims.Email, claims.Username)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate access token"})
		context.Abort()
		return
	}

	// Returns the new token
	context.JSON(http.StatusOK, gin.H{"access_token": newAccessToken})
}
