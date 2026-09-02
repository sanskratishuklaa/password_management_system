package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"

	"password_manager/config"
	"password_manager/internal/handler"
	"password_manager/internal/repository"
	"password_manager/internal/service"
)

func main() {

	// Connect to PostgreSQL database
	db, err := config.ConnectDatabase()
	if err != nil {
		log.Fatal(err)
	}

	// Close database connection when application stops
	defer db.Close(context.Background())

	log.Println("Database Connected Successfully")

	// Create repository
	userRepository := repository.NewUserRepository(db)

	// Create service
	userService := service.NewUserService(userRepository)

	// Create authentication handler
	authHandler := handler.NewAuthHandler(userService)

	// Create Gin router
	router := gin.Default()

	// Health route
	router.GET("/api/v1/health", handler.Health)

	// Authentication route
	router.POST("/api/v1/auth/register", authHandler.Register)

	// Start server
	log.Println("Server running on http://localhost:8080")

	err = router.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}