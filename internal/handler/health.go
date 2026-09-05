package handler

import "github.com/gin-gonic/gin"

type HealthResponse struct {
	Status string `json:"status"`
}

func Health(c *gin.Context) {

	response := HealthResponse{
		Status: "ok",
	}

	c.JSON(200, response)
}
