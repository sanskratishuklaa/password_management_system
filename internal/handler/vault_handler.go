package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"password_manager/internal/model"
	"password_manager/internal/service"
)

type VaultHandler struct {
	vaultService *service.VaultService
}

func NewVaultHandler(
	vaultService *service.VaultService,
) *VaultHandler {
	return &VaultHandler{
		vaultService: vaultService,
	}
}

func (h *VaultHandler) CreateVaultItem(c *gin.Context) {
	var request model.CreateVaultRequest

	// Read JSON body
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	// Get authenticated user ID from JWT middleware
	userID, exists := c.Get("userID")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user ID not found",
		})
		return
	}

	userIDString, ok := userID.(string)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid user ID",
		})
		return
	}

	// Call service
	vaultItem, err := h.vaultService.CreateVaultItem(
		c.Request.Context(),
		userIDString,
		request,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, vaultItem)
}

func (h *VaultHandler) GetVaultItems(c *gin.Context) {

	userID, exists := c.Get("userID")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user ID not found",
		})
		return
	}

	userIDString, ok := userID.(string)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid user ID",
		})
		return
	}

	vaultItems, err := h.vaultService.GetVaultItems(
		c.Request.Context(),
		userIDString,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": vaultItems,
	})
}
