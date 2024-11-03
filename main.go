package main

import "auth-service/database"

func main() {
	// Initialize Database
	database.Connect() 
	database.Migrate()
}