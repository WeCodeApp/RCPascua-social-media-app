package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go-azure/middleware"
	"go-azure/models"
	"go-azure/services"
	"go-azure/utils"
)

// TaskController handles task endpoints
type TaskController struct {
	taskService *services.TaskService
	authMiddleware *middleware.AuthMiddleware
	logger      *logrus.Logger
}

// NewTaskController creates a new TaskController
func NewTaskController(taskService *services.TaskService, authMiddleware *middleware.AuthMiddleware) *TaskController {
	return &TaskController{
		taskService: taskService,
		authMiddleware: authMiddleware,
		logger:      utils.GetLogger(),
	}
}

// RegisterRoutes registers the routes for the TaskController
func (c *TaskController) RegisterRoutes(router *gin.Engine) {
	tasks := router.Group("/tasks")
	tasks.Use(c.authMiddleware.RequireAuth())
	{
		tasks.GET("", c.GetAllTasks)
		tasks.GET("/:id", c.GetTaskByID)
		tasks.POST("", c.CreateTask)
		tasks.PUT("/:id", c.UpdateTask)
		tasks.DELETE("/:id", c.DeleteTask)
	}
}

// GetAllTasks returns all tasks for the authenticated user
func (c *TaskController) GetAllTasks(ctx *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID := ctx.GetString("user_id")

	// Get tasks
	tasks := c.taskService.GetAllTasks(userID)

	ctx.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

// GetTaskByID returns a task by ID
func (c *TaskController) GetTaskByID(ctx *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID := ctx.GetString("user_id")

	// Get task ID from URL
	taskID := ctx.Param("id")

	// Get task
	task, err := c.taskService.GetTaskByID(taskID, userID)
	if err != nil {
		c.logger.WithError(err).Error("Failed to get task")
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"task": task})
}

// CreateTask creates a new task
func (c *TaskController) CreateTask(ctx *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID := ctx.GetString("user_id")

	// Parse request body
	var task models.Task
	if err := ctx.ShouldBindJSON(&task); err != nil {
		c.logger.WithError(err).Error("Failed to parse request body")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create task
	createdTask, err := c.taskService.CreateTask(&task, userID)
	if err != nil {
		c.logger.WithError(err).Error("Failed to create task")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"task": createdTask})
}

// UpdateTask updates an existing task
func (c *TaskController) UpdateTask(ctx *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID := ctx.GetString("user_id")

	// Get task ID from URL
	taskID := ctx.Param("id")

	// Parse request body
	var task models.Task
	if err := ctx.ShouldBindJSON(&task); err != nil {
		c.logger.WithError(err).Error("Failed to parse request body")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update task
	updatedTask, err := c.taskService.UpdateTask(taskID, &task, userID)
	if err != nil {
		c.logger.WithError(err).Error("Failed to update task")
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"task": updatedTask})
}

// DeleteTask deletes a task
func (c *TaskController) DeleteTask(ctx *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID := ctx.GetString("user_id")

	// Get task ID from URL
	taskID := ctx.Param("id")

	// Delete task
	err := c.taskService.DeleteTask(taskID, userID)
	if err != nil {
		c.logger.WithError(err).Error("Failed to delete task")
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}