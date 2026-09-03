package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"password_manager/config"
	"password_manager/internal/handler"
	"password_manager/internal/middleware"
	"password_manager/internal/repository"
	"password_manager/internal/service"
)

func main() {

	// Load environment variables from .env
	config.LoadEnv()

	// Connect to PostgreSQL
	db, err := config.ConnectDatabase()
	if err != nil {
		log.Fatal(err)
	}

	// Close database connection when application stops
	defer db.Close(context.Background())

	log.Println("Database Connected Successfully")

	// Get JWT secret
	jwtSecret := os.Getenv("JWT_SECRET")

	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is not configured")
	}

	// Create repository
	userRepository := repository.NewUserRepository(db)

	// Create service
	userService := service.NewUserService(
		userRepository,
		jwtSecret,
	)

	// Create authentication handler
	authHandler := handler.NewAuthHandler(userService)

	// Create Gin router
	router := gin.Default()

	// Public routes
	router.GET(
		"/api/v1/health",
		handler.Health,
	)

	router.POST(
		"/api/v1/auth/register",
		authHandler.Register,
	)

	router.POST(
		"/api/v1/auth/login",
		authHandler.Login,
	)

	// JWT authentication middleware
	authMiddleware := middleware.AuthMiddleware(jwtSecret)

	// Protected routes
	protected := router.Group("/api/v1")
	protected.Use(authMiddleware)

	protected.GET(
		"/protected",
		handler.Protected,
	)

	// Start server
	log.Println("Server running on http://localhost:8080")

	err = router.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
