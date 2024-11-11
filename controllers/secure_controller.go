// Protected endpoint.
package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Ping(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{"super secret message": "pong!"})
}