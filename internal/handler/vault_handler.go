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
	userID, exists := c.Get("user_id")

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

	userID, exists := c.Get("user_id")

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

func (h *VaultHandler) GetVaultItemByID(c *gin.Context) {

	itemID := c.Param("id")

	if itemID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "vault item ID is required",
		})
		return
	}

	userID, exists := c.Get("user_id")

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

	vaultItem, err := h.vaultService.GetVaultItemByID(
		c.Request.Context(),
		itemID,
		userIDString,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "vault item not found",
		})
		return
	}

	c.JSON(http.StatusOK, vaultItem)
}

func (h *VaultHandler) UpdateVaultItem(c *gin.Context) {

	itemID := c.Param("id")

	if itemID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "vault item ID is required",
		})
		return
	}

	userID, exists := c.Get("user_id")

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

	var request model.UpdateVaultItemRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	vaultItem, err := h.vaultService.UpdateVaultItem(
		c.Request.Context(),
		itemID,
		userIDString,
		request.Title,
		request.Username,
		request.Password,
		request.Website,
		request.Notes,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, vaultItem)
}
func (h *VaultHandler) DeleteVaultItem(c *gin.Context) {

	itemID := c.Param("id")

	if itemID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "vault item ID is required",
		})
		return
	}

	userID, exists := c.Get("user_id")

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

	err := h.vaultService.DeleteVaultItem(
		c.Request.Context(),
		itemID,
		userIDString,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "vault item not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "vault item deleted successfully",
	})
}
