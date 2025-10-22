package handlers

import (
	"net/http"
	"twitter-bookmarks-api/database"
	"twitter-bookmarks-api/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ExportBookmarks(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	bookmarks, err := database.GetAllBookmarksByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to fetch bookmarks"})
		return
	}

	c.Header("Content-Disposition", "attachment; filename=bookmarks.json")
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, gin.H{"bookmarks": bookmarks})
}

func ExportCategory(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	categoryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid category ID"})
		return
	}

	category, err := database.GetCategoryByID(c.Request.Context(), categoryID, userID)
	if err != nil || category == nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Category not found"})
		return
	}

	bookmarks, err := database.GetBookmarksByCategoryID(c.Request.Context(), categoryID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to fetch bookmarks"})
		return
	}

	c.Header("Content-Disposition", "attachment; filename="+category.Name+".json")
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, gin.H{
		"category":  category,
		"bookmarks": bookmarks,
	})
}
