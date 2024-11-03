package controllers

import (
	"auth-service/database"
	"auth-service/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterUser(context *gin.Context) {
	var user models.User

	// Checks if the sent json structure matches the user structure.
	if err := context.ShouldBindJSON(&user); err != nil {
		// Error handling
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	// Hashes the password.
	if err := user.HashPassword(user.Password); err != nil {
		// Error handling
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	// Once hashed, the user data is stored into the DB using the GORM global instance.
	record := database.Instance.Create(&user) 
	if record.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": record.Error.Error()})
    	context.Abort()
    	return
	}

	// Returns status code 200 along with the created user data
	context.JSON(http.StatusCreated, gin.H{"userId": user.ID, "email": user.Email, "username": user.Username})
}