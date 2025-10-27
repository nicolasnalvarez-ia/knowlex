package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID `json:"id"`
	TwitterID      string    `json:"twitter_id"`
	Username       string    `json:"username"`
	DisplayName    string    `json:"display_name"`
	ProfileImage   string    `json:"profile_image"`
	AutoCategorize bool      `json:"auto_categorize"`
	CreatedAt      time.Time `json:"created_at"`
}

type Bookmark struct {
	ID                uuid.UUID  `json:"id"`
	UserID            uuid.UUID  `json:"user_id"`
	TweetID           string     `json:"tweet_id"`
	TweetText         string     `json:"tweet_text"`
	AuthorUsername    string     `json:"author_username"`
	AuthorDisplayName string     `json:"author_display_name"`
	TweetURL          string     `json:"tweet_url"`
	MediaURLs         []string   `json:"media_urls"`
	BookmarkedAt      time.Time  `json:"bookmarked_at"`
	CreatedAt         time.Time  `json:"created_at"`
	Categories        []Category `json:"categories,omitempty"`
}

type Category struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Name      string    `json:"name"`
	Color     string    `json:"color"`
	Icon      string    `json:"icon"`
	CreatedAt time.Time `json:"created_at"`
	Count     int       `json:"count,omitempty"`
}

type BookmarkImportItem struct {
	TweetID           string   `json:"tweet_id"`
	TweetText         string   `json:"tweet_text"`
	AuthorUsername    string   `json:"author_username"`
	AuthorDisplayName string   `json:"author_display_name"`
	TweetURL          string   `json:"tweet_url"`
	MediaURLs         []string `json:"media_urls"`
	BookmarkedAt      string   `json:"bookmarked_at"`
}

type BookmarkImport struct {
	Bookmarks []BookmarkImportItem `json:"bookmarks"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type ImportResponse struct {
	Message         string `json:"message"`
	ImportedCount   int    `json:"imported_count"`
	DuplicateCount  int    `json:"duplicate_count"`
	AutoCategorized int    `json:"auto_categorized,omitempty"`
}

type PaginationParams struct {
	Page     int
	PageSize int
	Offset   int
}

type BookmarksResponse struct {
	Bookmarks  []Bookmark `json:"bookmarks"`
	Total      int        `json:"total"`
	Page       int        `json:"page"`
	PageSize   int        `json:"page_size"`
	TotalPages int        `json:"total_pages"`
}

type CreateCategoryRequest struct {
	Name  string `json:"name" binding:"required"`
	Color string `json:"color"`
	Icon  string `json:"icon"`
}

type UpdateCategoryRequest struct {
	Name  string `json:"name"`
	Color string `json:"color"`
	Icon  string `json:"icon"`
}

type AssignCategoryRequest struct {
	CategoryID string `json:"category_id" binding:"required"`
}

type UpdatePreferencesRequest struct {
	AutoCategorize *bool `json:"auto_categorize"`
}
