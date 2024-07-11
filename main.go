package main

import (
	"example/users/database"
	"example/users/handlers"
	"example/users/middleware"
	"example/users/models"
	"log"

	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
)

func main() {

	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect to the database
	database.Connect()

	// Auto migrate the database schema
	database.DB.AutoMigrate(&models.User{})

	r := gin.Default()

	r.GET("/users", middleware.RequireAuth, handlers.GetUsers)
	r.POST("/users", middleware.RequireAuth, handlers.AddUser)
	r.POST("/login", handlers.Login)

	r.Run(":8080")
}
