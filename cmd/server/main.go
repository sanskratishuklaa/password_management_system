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

	db, err := config.ConnectDatabase()
	if err != nil {
		log.Fatal(err)
	}


	defer db.Close(context.Background())

	log.Println("Database Connected Successfully")


	userRepository := repository.NewUserRepository(db)
	
	userService := service.NewUserService(userRepository)
	authHandler := handler.NewAuthHandler(userService)

	
	router := gin.Default()

	
	router.GET("/api/v1/health", handler.Health)


	router.POST("/api/v1/auth/register", authHandler.Register)
	router.POST("/api/v1/auth/login", authHandler.Login)

	
	log.Println("Server running on http://localhost:8080")

	err = router.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}