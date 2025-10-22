package handlers

import (
	"net/http"
	"twitter-bookmarks-api/database"
	"twitter-bookmarks-api/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func DeleteAccount(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	err := database.DeleteUserAndAllData(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to delete account"})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Message: "Account deleted successfully"})
}
