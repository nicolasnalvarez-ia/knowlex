package handlers

import (
	"net/http"
	"twitter-bookmarks-api/database"
	"twitter-bookmarks-api/models"
	"twitter-bookmarks-api/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AutoCategorizeBookmarks uses AI to categorize uncategorized bookmarks
func AutoCategorizeBookmarks(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	// Get user's existing categories
	// Get uncategorized bookmarks (bookmarks with no categories)
	bookmarks, err := database.GetUncategorizedBookmarks(c.Request.Context(), userID, 50) // Limit to 50 at a time
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to fetch bookmarks"})
		return
	}

	if len(bookmarks) == 0 {
		c.JSON(http.StatusOK, models.SuccessResponse{
			Message: "No uncategorized bookmarks found",
			Data:    map[string]int{"categorized": 0},
		})
		return
	}

	categorizedCount, newCategoriesCount, err := services.CategorizeBookmarksForUser(c.Request.Context(), userID, bookmarks)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to categorize bookmarks"})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Auto-categorization completed",
		Data: map[string]int{
			"categorized":     categorizedCount,
			"new_categories":  newCategoriesCount,
			"total_processed": len(bookmarks),
		},
	})
}

// CategorizeBookmark categorizes a single bookmark with AI
func CategorizeBookmark(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	bookmarkID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid bookmark ID"})
		return
	}

	// Get the bookmark
	bookmark, err := database.GetBookmarkByID(c.Request.Context(), bookmarkID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Bookmark not found"})
		return
	}

	// Use services helper to categorize this single bookmark
	categorized, _, err := services.CategorizeBookmarksForUser(c.Request.Context(), userID, []models.Bookmark{*bookmark})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "AI categorization failed"})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Categorization completed",
		Data: map[string]int{
			"categorized": categorized,
		},
	})
}
