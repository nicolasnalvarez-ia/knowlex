package handlers

import (
	"net/http"
	"twitter-bookmarks-api/ai"
	"twitter-bookmarks-api/database"
	"twitter-bookmarks-api/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AutoCategorizeBookmarks uses AI to categorize uncategorized bookmarks
func AutoCategorizeBookmarks(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	// Get user's existing categories
	categories, err := database.GetCategoriesByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to fetch categories"})
		return
	}

	categoryNames := make([]string, len(categories))
	categoryMap := make(map[string]uuid.UUID)
	for i, cat := range categories {
		categoryNames[i] = cat.Name
		categoryMap[cat.Name] = cat.ID
	}

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

	categorizedCount := 0
	newCategoriesCount := 0

	// Process each bookmark
	for _, bookmark := range bookmarks {
		suggestedCategories, err := ai.CategorizeBookmark(c.Request.Context(), bookmark.TweetText, categoryNames)
		if err != nil {
			continue // Skip on error
		}

		// Assign categories
		for _, catName := range suggestedCategories {
			// Check if category exists
			catID, exists := categoryMap[catName]
			if !exists {
				// Create new category
				newCategory := &models.Category{
					UserID: userID,
					Name:   catName,
					Color:  getColorForCategory(catName),
					Icon:   getIconForCategory(catName),
				}
				
				err = database.CreateCategory(c.Request.Context(), newCategory)
				if err != nil {
					continue
				}
				
				catID = newCategory.ID
				categoryMap[catName] = catID
				categoryNames = append(categoryNames, catName)
				newCategoriesCount++
			}

			// Assign bookmark to category
			err = database.AssignBookmarkToCategory(c.Request.Context(), bookmark.ID, catID, userID)
			if err == nil {
				categorizedCount++
			}
		}
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Auto-categorization completed",
		Data: map[string]int{
			"categorized":      categorizedCount,
			"new_categories":   newCategoriesCount,
			"total_processed":  len(bookmarks),
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

	// Get user's existing categories
	categories, err := database.GetCategoriesByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to fetch categories"})
		return
	}

	categoryNames := make([]string, len(categories))
	for i, cat := range categories {
		categoryNames[i] = cat.Name
	}

	// Get AI suggestions
	suggestedCategories, err := ai.CategorizeBookmark(c.Request.Context(), bookmark.TweetText, categoryNames)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "AI categorization failed"})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Categorization suggestions generated",
		Data:    map[string][]string{"suggested_categories": suggestedCategories},
	})
}

// Helper functions for default colors and icons based on category name
func getColorForCategory(name string) string {
	colorMap := map[string]string{
		"Tech":      "#3B82F6", // Blue
		"AI":        "#8B5CF6", // Purple
		"Design":    "#EC4899", // Pink
		"Coding":    "#10B981", // Green
		"Business":  "#F59E0B", // Amber
		"Marketing": "#EF4444", // Red
		"News":      "#6366F1", // Indigo
		"Tutorial":  "#14B8A6", // Teal
		"Meme":      "#F97316", // Orange
		"Thread":    "#A855F7", // Purple
	}

	if color, ok := colorMap[name]; ok {
		return color
	}
	return "#6B7280" // Default gray
}

func getIconForCategory(name string) string {
	iconMap := map[string]string{
		"Tech":      "💻",
		"AI":        "🤖",
		"Design":    "🎨",
		"Coding":    "⚡",
		"Business":  "💼",
		"Marketing": "📈",
		"News":      "📰",
		"Tutorial":  "📚",
		"Meme":      "😂",
		"Thread":    "🧵",
	}

	if icon, ok := iconMap[name]; ok {
		return icon
	}
	return "📁" // Default folder
}

