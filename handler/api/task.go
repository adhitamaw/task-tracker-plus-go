package api

import (
	"a21hc3NpZ25tZW50/model"
	"a21hc3NpZ25tZW50/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TaskAPI interface {
	AddTask(c *gin.Context)
	UpdateTask(c *gin.Context)
	DeleteTask(c *gin.Context)
	GetTaskByID(c *gin.Context)
	GetTaskList(c *gin.Context)
	GetTaskListByCategory(c *gin.Context)
}

type taskAPI struct {
	taskService service.TaskService
}

func NewTaskAPI(taskRepo service.TaskService) *taskAPI {
	return &taskAPI{taskRepo}
}

func (t *taskAPI) AddTask(c *gin.Context) {
	var newTask model.Task
	if err := c.ShouldBindJSON(&newTask); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error()})
		return
	}

	// Validasi data task
	if newTask.Title == "" {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Title cannot be empty"})
		return
	}

	if newTask.Deadline == "" {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Deadline cannot be empty"})
		return
	}

	if newTask.Status == "" {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Status cannot be empty"})
		return
	}

	if newTask.CategoryID <= 0 {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid category ID"})
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

	// Set the user ID from the authenticated user
	newTask.UserID = userIDInt

	err := t.taskService.Store(&newTask)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{Message: "add task success"})
}

func (t *taskAPI) UpdateTask(c *gin.Context) {
	taskID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid task ID"})
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

	// First, check if task exists and belongs to user
	existingTask, err := t.taskService.GetByID(taskID)
	if err != nil {
		c.JSON(http.StatusNotFound, model.ErrorResponse{Error: "Task not found"})
		return
	}

	if existingTask.UserID != userIDInt {
		c.JSON(http.StatusForbidden, model.ErrorResponse{Error: "Access denied: task belongs to different user"})
		return
	}

	var updatedTask model.Task
	if err := c.ShouldBindJSON(&updatedTask); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error()})
		return
	}

	updatedTask.ID = taskID
	updatedTask.UserID = userIDInt // Ensure user ID remains the same
	err = t.taskService.Update(taskID, &updatedTask)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{Message: "update task success"})
}

func (t *taskAPI) DeleteTask(c *gin.Context) {
	taskID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid task ID"})
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

	// First, check if task exists and belongs to user
	existingTask, err := t.taskService.GetByID(taskID)
	if err != nil {
		c.JSON(http.StatusNotFound, model.ErrorResponse{Error: "Task not found"})
		return
	}

	if existingTask.UserID != userIDInt {
		c.JSON(http.StatusForbidden, model.ErrorResponse{Error: "Access denied: task belongs to different user"})
		return
	}

	err = t.taskService.Delete(taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{Message: "delete task success"})
}

func (t *taskAPI) GetTaskByID(c *gin.Context) {
	taskID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid task ID"})
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

	task, err := t.taskService.GetByID(taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	// Validate task ownership
	if task.UserID != userIDInt {
		c.JSON(http.StatusForbidden, model.ErrorResponse{Error: "Access denied: task belongs to different user"})
		return
	}

	c.JSON(http.StatusOK, task)
}

func (t *taskAPI) GetTaskList(c *gin.Context) {
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

	tasks, err := t.taskService.GetList(userIDInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func (t *taskAPI) GetTaskListByCategory(c *gin.Context) {
	categoryID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid category ID"})
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

	taskCategories, err := t.taskService.GetTaskCategoryByUser(categoryID, userIDInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, taskCategories)
}
