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
	// Connect to database
	db, err := config.ConnectDatabase()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close(context.Background())

	log.Println("Database Connected Successfully")

	// Get JWT secret
	jwtSecret := os.Getenv("JWT_SECRET")

	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is not set")
	}

	userRepository := repository.NewUserRepository(db)

	vaultRepository := repository.NewVaultRepository(db)

	userService := service.NewUserService(
		userRepository,
		jwtSecret,
	)

	vaultService := service.NewVaultService(
		vaultRepository,
	)

	authHandler := handler.NewAuthHandler(
		userService,
	)

	vaultHandler := handler.NewVaultHandler(
		vaultService,
	)

	router := gin.Default()

	// Health
	router.GET(
		"/api/v1/health",
		handler.Health,
	)

	// Authentication
	router.POST(
		"/api/v1/auth/register",
		authHandler.Register,
	)

	router.POST(
		"/api/v1/auth/login",
		authHandler.Login,
	)

	protected := router.Group("/api/v1")

	protected.Use(
		middleware.AuthMiddleware(jwtSecret),
	)

	protected.POST(
		"/vault",
		vaultHandler.CreateVaultItem,
	)

	protected.GET(
		"/vault",
		vaultHandler.GetVaultItems,
	)
	protected.GET(
		"/vault/:id",
		vaultHandler.GetVaultItemByID,
	)
	protected.PUT(
		"/vault/:id",
		vaultHandler.UpdateVaultItem,
	)
	protected.DELETE(
		"/vault/:id",
		vaultHandler.DeleteVaultItem,
	)

	log.Println("Server running on http://localhost:8080")

	err = router.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
