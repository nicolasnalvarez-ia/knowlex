package services

import (
	"context"
	"strings"
	"twitter-bookmarks-api/ai"
	"twitter-bookmarks-api/database"
	"twitter-bookmarks-api/models"

	"github.com/google/uuid"
)

// CategorizeBookmarksForUser applies AI categorization to the provided bookmarks.
// Returns the number of bookmark-category assignments made and new categories created.
func CategorizeBookmarksForUser(ctx context.Context, userID uuid.UUID, bookmarks []models.Bookmark) (int, int, error) {
	if len(bookmarks) == 0 {
		return 0, 0, nil
	}

	categories, err := database.GetCategoriesByUserID(ctx, userID)
	if err != nil {
		return 0, 0, err
	}

	categoryNames := make([]string, len(categories))
	categoryMap := make(map[string]uuid.UUID, len(categories))

	for i, cat := range categories {
		categoryNames[i] = cat.Name
		categoryMap[strings.ToLower(cat.Name)] = cat.ID
	}

	categorizedCount := 0
	newCategoriesCount := 0

	for _, bookmark := range bookmarks {
		if strings.TrimSpace(bookmark.TweetText) == "" {
			continue
		}

		suggestedCategories, err := ai.CategorizeBookmark(ctx, bookmark.TweetText, categoryNames)
		if err != nil {
			continue
		}

		for _, catName := range suggestedCategories {
			if catName == "" {
				continue
			}

			lookupKey := strings.ToLower(catName)
			catID, exists := categoryMap[lookupKey]

			if !exists {
				// Create new category
				newCategory := &models.Category{
					UserID: userID,
					Name:   catName,
					Color:  getColorForCategory(catName),
					Icon:   getIconForCategory(catName),
				}

				if err := database.CreateCategory(ctx, newCategory); err != nil {
					continue
				}

				catID = newCategory.ID
				categoryMap[lookupKey] = catID
				categoryNames = append(categoryNames, catName)
				newCategoriesCount++
			}

			if err := database.AssignBookmarkToCategory(ctx, bookmark.ID, catID, userID); err == nil {
				categorizedCount++
			}
		}
	}

	return categorizedCount, newCategoriesCount, nil
}

// Helper functions for default colors and icons based on category name
func getColorForCategory(name string) string {
	colorMap := map[string]string{
		"tech":      "#3B82F6", // Blue
		"ai":        "#8B5CF6", // Purple
		"design":    "#EC4899", // Pink
		"coding":    "#10B981", // Green
		"business":  "#F59E0B", // Amber
		"marketing": "#EF4444", // Red
		"news":      "#6366F1", // Indigo
		"tutorial":  "#14B8A6", // Teal
		"meme":      "#F97316", // Orange
		"thread":    "#A855F7", // Purple
	}

	if color, ok := colorMap[strings.ToLower(name)]; ok {
		return color
	}
	return "#6B7280" // Default gray
}

func getIconForCategory(name string) string {
	iconMap := map[string]string{
		"tech":      "💻",
		"ai":        "🤖",
		"design":    "🎨",
		"coding":    "⚡",
		"business":  "💼",
		"marketing": "📈",
		"news":      "📰",
		"tutorial":  "📚",
		"meme":      "😂",
		"thread":    "🧵",
	}

	if icon, ok := iconMap[strings.ToLower(name)]; ok {
		return icon
	}
	return "📁"
}
