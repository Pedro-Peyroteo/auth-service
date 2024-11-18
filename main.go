package main

import (
	"auth-service/controllers"
	"auth-service/database"
	"auth-service/middleware"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

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

	// Get allowed origins from .env file
	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	origins := strings.Split(allowedOrigins, ",")

	// Set up CORS middleware with configuration
	router.Use(cors.New(cors.Config{
		AllowOrigins:     origins, // Use origins from .env file
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true, // Allows cookies to be shared
		MaxAge:           12 * time.Hour,
	}))

	// Group everything under "/api"
	api := router.Group("/api")
	{
		// Routes endpoint to the correspondant controller
		api.POST("/token", controllers.GenerateToken)
		api.POST("/token/refresh", controllers.GenerateRefreshToken)
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
