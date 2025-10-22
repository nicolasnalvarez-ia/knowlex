package handlers

import (
	"net/http"
	"strconv"
	"time"
	"twitter-bookmarks-api/database"
	"twitter-bookmarks-api/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetBookmarks(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	categoryIDStr := c.Query("category_id")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	params := models.PaginationParams{
		Page:     page,
		PageSize: pageSize,
		Offset:   (page - 1) * pageSize,
	}

	var categoryID *uuid.UUID
	if categoryIDStr != "" {
		id, err := uuid.Parse(categoryIDStr)
		if err == nil {
			categoryID = &id
		}
	}

	response, err := database.GetBookmarksByUserID(c.Request.Context(), userID, params, categoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to fetch bookmarks"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func ImportBookmarks(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	var importData models.BookmarkImport
	if err := c.ShouldBindJSON(&importData); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid JSON format"})
		return
	}

	importedCount := 0
	duplicateCount := 0

	for _, item := range importData.Bookmarks {
		bookmarkedAt, err := time.Parse(time.RFC3339, item.BookmarkedAt)
		if err != nil {
			bookmarkedAt = time.Now()
		}

		bookmark := &models.Bookmark{
			UserID:            userID,
			TweetID:           item.TweetID,
			TweetText:         item.TweetText,
			AuthorUsername:    item.AuthorUsername,
			AuthorDisplayName: item.AuthorDisplayName,
			TweetURL:          item.TweetURL,
			MediaURLs:         item.MediaURLs,
			BookmarkedAt:      bookmarkedAt,
		}

		err = database.CreateBookmark(c.Request.Context(), bookmark)
		if err != nil {
			duplicateCount++
		} else {
			importedCount++
		}
	}

	c.JSON(http.StatusOK, models.ImportResponse{
		Message:        "Import completed",
		ImportedCount:  importedCount,
		DuplicateCount: duplicateCount,
	})
}

func DeleteBookmark(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	bookmarkID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid bookmark ID"})
		return
	}

	err = database.DeleteBookmark(c.Request.Context(), bookmarkID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Bookmark not found"})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Message: "Bookmark deleted successfully"})
}

func SearchBookmarks(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	query := c.Query("q")

	if query == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Search query is required"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	params := models.PaginationParams{
		Page:     page,
		PageSize: pageSize,
		Offset:   (page - 1) * pageSize,
	}

	response, err := database.SearchBookmarks(c.Request.Context(), userID, query, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to search bookmarks"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func AssignCategory(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	bookmarkID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid bookmark ID"})
		return
	}

	var req models.AssignCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid request body"})
		return
	}

	categoryID, err := uuid.Parse(req.CategoryID)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid category ID"})
		return
	}

	err = database.AssignBookmarkToCategory(c.Request.Context(), bookmarkID, categoryID, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Message: "Category assigned successfully"})
}

func RemoveCategory(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	bookmarkID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid bookmark ID"})
		return
	}

	categoryID, err := uuid.Parse(c.Param("categoryId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid category ID"})
		return
	}

	err = database.RemoveBookmarkFromCategory(c.Request.Context(), bookmarkID, categoryID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to remove category"})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Message: "Category removed successfully"})
}
