package services

import (
	"errors"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"go-azure/models"
	"go-azure/utils"
	"gorm.io/gorm"
)

// TaskService handles task operations
type TaskService struct {
	db     *gorm.DB
	logger *logrus.Logger
}

// NewTaskService creates a new TaskService
func NewTaskService() *TaskService {
	return &TaskService{
		db:     utils.GetDB(),
		logger: utils.GetLogger(),
	}
}

// GetAllTasks returns all tasks for a user
func (s *TaskService) GetAllTasks(userID string) []*models.Task {
	var tasks []*models.Task

	result := s.db.Where("user_id = ?", userID).Find(&tasks)
	if result.Error != nil {
		s.logger.WithError(result.Error).Error("Failed to get tasks")
		return []*models.Task{}
	}

	return tasks
}

// GetTaskByID returns a task by ID
func (s *TaskService) GetTaskByID(taskID string, userID string) (*models.Task, error) {
	var task models.Task

	result := s.db.Where("id = ? AND user_id = ?", taskID, userID).First(&task)
	if result.Error != nil {
		s.logger.WithError(result.Error).Error("Failed to get task")
		return nil, errors.New("task not found")
	}

	return &task, nil
}

// CreateTask creates a new task
func (s *TaskService) CreateTask(task *models.Task, userID string) (*models.Task, error) {
	// Set task ID and user ID
	task.ID = uuid.New().String()
	task.UserID = userID

	// Create task in database
	result := s.db.Create(task)
	if result.Error != nil {
		s.logger.WithError(result.Error).Error("Failed to create task")
		return nil, errors.New("failed to create task")
	}

	s.logger.WithFields(logrus.Fields{
		"task_id": task.ID,
		"user_id": userID,
	}).Info("Task created")

	return task, nil
}

// UpdateTask updates an existing task
func (s *TaskService) UpdateTask(taskID string, updatedTask *models.Task, userID string) (*models.Task, error) {
	// Get existing task
	var existingTask models.Task
	result := s.db.Where("id = ? AND user_id = ?", taskID, userID).First(&existingTask)
	if result.Error != nil {
		s.logger.WithError(result.Error).Error("Failed to get task for update")
		return nil, errors.New("task not found")
	}

	// Update task fields
	existingTask.Title = updatedTask.Title
	existingTask.Description = updatedTask.Description
	existingTask.Completed = updatedTask.Completed

	// Save changes to database
	result = s.db.Save(&existingTask)
	if result.Error != nil {
		s.logger.WithError(result.Error).Error("Failed to update task")
		return nil, errors.New("failed to update task")
	}

	s.logger.WithFields(logrus.Fields{
		"task_id": taskID,
		"user_id": userID,
	}).Info("Task updated")

	return &existingTask, nil
}

// DeleteTask deletes a task
func (s *TaskService) DeleteTask(taskID string, userID string) error {
	// Check if task exists and belongs to user
	var task models.Task
	result := s.db.Where("id = ? AND user_id = ?", taskID, userID).First(&task)
	if result.Error != nil {
		s.logger.WithError(result.Error).Error("Failed to get task for deletion")
		return errors.New("task not found")
	}

	// Delete task
	result = s.db.Delete(&task)
	if result.Error != nil {
		s.logger.WithError(result.Error).Error("Failed to delete task")
		return errors.New("failed to delete task")
	}

	s.logger.WithFields(logrus.Fields{
		"task_id": taskID,
		"user_id": userID,
	}).Info("Task deleted")

	return nil
}
