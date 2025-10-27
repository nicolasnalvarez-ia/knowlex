package handlers

import (
	"net/http"
	"twitter-bookmarks-api/database"
	"twitter-bookmarks-api/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func UpdateUserPreferences(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	var req models.UpdatePreferencesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid request body"})
		return
	}

	if err := database.UpdateUserPreferences(c.Request.Context(), userID, req.AutoCategorize); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to update preferences"})
		return
	}

	user, err := database.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to fetch user"})
		return
	}

	c.JSON(http.StatusOK, user)
}
