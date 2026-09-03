package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Protected(c *gin.Context) {
	userID := c.GetString("user_id")

	c.JSON(http.StatusOK, gin.H{
		"message": "you accessed a protected route",
		"user_id": userID,
	})
}
