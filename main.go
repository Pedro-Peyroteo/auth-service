package main

import (
	"auth-service/controllers"
	"auth-service/database"
	"auth-service/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize Database
	database.Connect()
	database.Migrate()

	// Initialize Router
	router := initRouter()
	// Router port
	router.Run(":8080")
}

func initRouter() *gin.Engine {
	// Creates a new Gin Router instance
	router := gin.Default()

	// Group everything under "/api"
	api := router.Group("/api")
	{
		// Routes endpoint to the correspondant controller
		api.POST("/token", controllers.GenerateToken)
		api.POST("/user/register", controllers.RegisterUser)
		// Adds a protected route, using the middleware implemmented
		secured := api.Group("/secured").Use(middleware.Auth())
		{
			// Middleware secured endpoints
			secured.GET("/ping", controllers.Ping)
		}
	}

	return router
}
