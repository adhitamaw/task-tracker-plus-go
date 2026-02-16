package api

import (
	"a21hc3NpZ25tZW50/model"
	"a21hc3NpZ25tZW50/service"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type CategoryAPI interface {
	AddCategory(c *gin.Context)
	UpdateCategory(c *gin.Context)
	DeleteCategory(c *gin.Context)
	GetCategoryByID(c *gin.Context)
	GetCategoryList(c *gin.Context)
}

type categoryAPI struct {
	categoryService service.CategoryService
}

func NewCategoryAPI(categoryRepo service.CategoryService) *categoryAPI {
	return &categoryAPI{categoryRepo}
}

func (ct *categoryAPI) AddCategory(c *gin.Context) {
	var newCategory model.Category
	if err := c.ShouldBindJSON(&newCategory); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error()})
		return
	}

	// Get user ID from context (set by middleware)
	userID, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Unauthorized"})
		return
	}

	userIDInt, ok := userID.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Invalid user ID format"})
		return
	}

	// Set the user ID for the new category
	newCategory.UserID = userIDInt

	// Check if user already has a category with this name (including default categories)
	existingCategories, err := ct.categoryService.GetListByUser(userIDInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	// Check for duplicate names (case-insensitive)
	newCategoryNameLower := strings.ToLower(strings.TrimSpace(newCategory.Name))
	for _, category := range existingCategories {
		if strings.ToLower(strings.TrimSpace(category.Name)) == newCategoryNameLower {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Category with this name already exists"})
			return
		}
	}

	err = ct.categoryService.Store(&newCategory)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{Message: "add category success"})
}

func (ct *categoryAPI) UpdateCategory(c *gin.Context) {
	categoryID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid Category ID"})
		return
	}

	// Get user ID from context (set by middleware)
	userID, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Unauthorized"})
		return
	}

	userIDInt, ok := userID.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Invalid user ID format"})
		return
	}

	// Check if category exists and belongs to user
	existingCategory, err := ct.categoryService.GetByID(categoryID)
	if err != nil {
		c.JSON(http.StatusNotFound, model.ErrorResponse{Error: "Category not found"})
		return
	}

	// Only allow users to update their own categories (not system categories with UserID=0)
	if existingCategory.UserID != userIDInt {
		c.JSON(http.StatusForbidden, model.ErrorResponse{Error: "Access denied: category belongs to different user or is a system category"})
		return
	}

	var updatedCategory model.Category
	if err := c.ShouldBindJSON(&updatedCategory); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error()})
		return
	}

	updatedCategory.ID = categoryID
	updatedCategory.UserID = userIDInt // Ensure user ID remains the same
	err = ct.categoryService.Update(categoryID, updatedCategory)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{Message: "category update success"})
}

func (ct *categoryAPI) DeleteCategory(c *gin.Context) {
	categoryID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid Category ID"})
		return
	}

	// Get user ID from context (set by middleware)
	userID, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Unauthorized"})
		return
	}

	userIDInt, ok := userID.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Invalid user ID format"})
		return
	}

	// Check if category exists and belongs to user
	existingCategory, err := ct.categoryService.GetByID(categoryID)
	if err != nil {
		c.JSON(http.StatusNotFound, model.ErrorResponse{Error: "Category not found"})
		return
	}

	// Only allow users to delete their own categories (not system categories with UserID=0)
	if existingCategory.UserID != userIDInt {
		c.JSON(http.StatusForbidden, model.ErrorResponse{Error: "Access denied: category belongs to different user or is a system category"})
		return
	}

	err = ct.categoryService.Delete(categoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{Message: "category delete success"})
}

func (ct *categoryAPI) GetCategoryByID(c *gin.Context) {
	categoryID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid Category ID"})
		return
	}

	userID, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Unauthorized"})
		return
	}

	userIDInt, ok := userID.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Invalid user ID format"})
		return
	}

	category, err := ct.categoryService.GetByID(categoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	if category.UserID != userIDInt {
		c.JSON(http.StatusForbidden, model.ErrorResponse{Error: "Access denied: category belongs to different user or is a system category"})
		return
	}

	c.JSON(http.StatusOK, category)
}

func (ct *categoryAPI) GetCategoryList(c *gin.Context) {
	// Get user ID from context (set by middleware)
	userID, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Unauthorized"})
		return
	}

	userIDInt, ok := userID.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Invalid user ID format"})
		return
	}

	// Get categories for this specific user (includes system categories + user's own)
	categories, err := ct.categoryService.GetListByUser(userIDInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, categories)
}
