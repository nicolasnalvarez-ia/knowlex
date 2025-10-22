package handlers

import (
	"net/http"
	"twitter-bookmarks-api/database"
	"twitter-bookmarks-api/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetCategories(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	categories, err := database.GetCategoriesByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to fetch categories"})
		return
	}

	if categories == nil {
		categories = []models.Category{}
	}

	c.JSON(http.StatusOK, categories)
}

func CreateCategory(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	var req models.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid request body"})
		return
	}

	if req.Color == "" {
		req.Color = "#3B82F6"
	}
	if req.Icon == "" {
		req.Icon = "folder"
	}

	category := &models.Category{
		UserID: userID,
		Name:   req.Name,
		Color:  req.Color,
		Icon:   req.Icon,
	}

	err := database.CreateCategory(c.Request.Context(), category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to create category"})
		return
	}

	c.JSON(http.StatusCreated, category)
}

func UpdateCategory(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	categoryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid category ID"})
		return
	}

	var req models.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid request body"})
		return
	}

	err = database.UpdateCategory(c.Request.Context(), categoryID, userID, req.Name, req.Color, req.Icon)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Category not found"})
		return
	}

	category, _ := database.GetCategoryByID(c.Request.Context(), categoryID, userID)
	c.JSON(http.StatusOK, category)
}

func DeleteCategory(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	categoryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid category ID"})
		return
	}

	err = database.DeleteCategory(c.Request.Context(), categoryID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Category not found"})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Message: "Category deleted successfully"})
}
