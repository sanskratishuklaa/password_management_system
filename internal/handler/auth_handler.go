package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"password_manager/internal/model"
	"password_manager/internal/service"
)

type AuthHandler struct {
	userService *service.UserService
}

func NewAuthHandler(userService *service.UserService) *AuthHandler {
	return &AuthHandler{
		userService: userService,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {

	var request model.RegisterRequest

	// Read JSON body
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	// Call service
	user, err := h.userService.RegisterUser(
		c.Request.Context(),
		request,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Success response
	c.JSON(http.StatusCreated, user)
}