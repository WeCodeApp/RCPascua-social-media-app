package controllers

import (
	"net/http"
	"strconv"

	"go-azure/middleware"
	"go-azure/models"
	"go-azure/services"
	"go-azure/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// SocialMediaController handles social media endpoints
type SocialMediaController struct {
	socialmediaService *services.SocialMediaService
	authMiddleware     *middleware.AuthMiddleware
	logger             *logrus.Logger
}

// NewSocialMediaController creates a new SocialMediaController
func NewSocialMediaController(socialmediaService *services.SocialMediaService, authMiddleware *middleware.AuthMiddleware) *SocialMediaController {
	return &SocialMediaController{
		socialmediaService: socialmediaService,
		authMiddleware:     authMiddleware,
		logger:             utils.GetLogger(),
	}
}

// RegisterRoutes registers the routes for the SocialMediaController
func (c *SocialMediaController) RegisterRoutes(router *gin.Engine) {
	posts := router.Group("/posts")
	posts.Use(c.authMiddleware.RequireAuth())
	{
		posts.GET("/page/:page_num/:page_limit", c.GetAllSocialMediaPosts)
		posts.GET("/page/:page_num/:page_limit/:sort_by/:sort_order", c.GetAllSocialMediaPosts)
		posts.GET("", c.GetAllSocialMediaPosts)
		posts.GET("/:post_id/", c.GetSocialMediaPostByPostID)
		posts.GET("/user/:user_id", c.GetAllSocialMediaPostByUserID)
		posts.GET("/:post_id/user", c.GetSocialMediaPostByPostAndUserID)
		posts.POST("", c.CreateSocialMediaPost)
		posts.PUT("/:post_id", c.UpdateSocialMediaPost)
		posts.DELETE("/:post_id", c.DeleteSocialMediaPost)
	}
}

// QuerySocialMediaPost queries social media posts with filters
// @Summary Query social media posts
// @Description Query social media posts with pagination, sorting, and filtering options
// @Tags SocialMedia
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of posts per page" default(10)
// @Param colname query string false "Column name to filter by"
// @Param searchtext query string false "Search text for filtering"
// @Param sort_by query string false "Field to sort by" default(created_at)
// @Param sort_order query string false "Sort order (asc or desc)" default(desc)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
func (c *SocialMediaController) QuerySocialMediaPost(ctx *gin.Context) {
	// USAGE: http://localhost:8080/posts?page=1&limit=10&colname=post_text&searchtext=test&create_at&desc
	// Extract query parameters
	pageStr := ctx.DefaultQuery("page", "1")
	limitStr := ctx.DefaultQuery("limit", "10")
	colName := ctx.Query("colname")
	searchText := ctx.Query("searchtext")
	sortBy := ctx.DefaultQuery("sort_by", "created_at")
	sortOrder := ctx.DefaultQuery("sort_order", "desc")

	// Convert page and limit to integers
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		c.logger.WithError(err).Error("Invalid page parameter")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page parameter"})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.logger.WithError(err).Error("Invalid limit parameter")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
		return
	}

	// Call the service to query posts
	posts, err := c.socialmediaService.QuerySocialMediaPost(page, limit, colName, searchText, sortBy, sortOrder)
	if err != nil {
		c.logger.WithError(err).Error("Failed to query social media posts")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the posts as JSON
	ctx.JSON(http.StatusOK, gin.H{"posts": posts})
}

// GetAllSocialMediaPosts retrieves social media posts with pagination and sorting
// @Summary Retrieve social media posts
// @Description Fetches a list of social media posts with pagination and sorting options
// @Tags SocialMedia
// @Accept json
// @Produce json
// @Param page_num path int false "Page number" default(1)
// @Param page_limit path int false "Number of posts per page" default(10)
// @Param sort_by path string false "Field to sort by" default(created_at)
// @Param sort_order path string false "Sort order (asc or desc)" default(desc)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
func (c *SocialMediaController) GetAllSocialMediaPosts(ctx *gin.Context) {
	// Extract path parameters
	pageNum := ctx.Param("page_num")
	if pageNum == "" {
		pageNum = "1" // Default page number
	}

	pageLimit := ctx.Param("page_limit")
	if pageLimit == "" {
		pageLimit = "10" // Default page limit
	}

	sortBy := ctx.Param("sort_by")
	if sortBy == "" {
		sortBy = "created_at" // Default sort by field
	}

	sortOrder := ctx.Param("sort_order")
	if sortOrder == "" {
		sortOrder = "desc" // Default sort order
	}

	// Log the extracted parameters for debugging
	c.logger.Infof("Page Number: %s, Page Limit: %s, Sort By: %s, Sort Order: %s", pageNum, pageLimit, sortBy, sortOrder)

	// Convert page and limit to integers
	page, err := strconv.Atoi(pageNum)
	if err != nil {
		c.logger.WithError(err).Error("Invalid page number parameter")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number parameter"})
		return
	}

	limit, err := strconv.Atoi(pageLimit)
	if err != nil {
		c.logger.WithError(err).Error("Invalid page limit parameter")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page limit parameter"})
		return
	}

	// Call the service to get posts
	response, err := c.socialmediaService.GetAllSocialMediaPosts(page, limit, sortBy, sortOrder)
	if err != nil {
		c.logger.WithError(err).Error("Failed to get social media posts")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the response as JSON
	ctx.JSON(http.StatusOK, response)
}

// GetSocialMediaPostByPostID retrieves a social media post by its ID
// @Summary Retrieve a social media post by ID
// @Description Fetches a single social media post by its ID
// @Tags SocialMedia
// @Accept json
// @Produce json
// @Param post_id path string true "Post ID"
// @Success 200 {object} models.SocialMediaPost
// @Failure 404 {object} gin.H
func (c *SocialMediaController) GetSocialMediaPostByPostID(ctx *gin.Context) {
	// Get postID from URL
	postID := ctx.Param("post_id")

	// Get posts details by postID
	post, err := c.socialmediaService.GetSocialMediaPostByPostID(postID)
	if err != nil {
		c.logger.WithError(err).Error("Failed to get task")
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if post.PostID == "" {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "Post Not Found"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"post": post})
}

// GetAllSocialMediaPostByUserID retrieves all social media posts for a user
// @Summary Retrieve all social media posts for a user
// @Description Fetches all social media posts created by a specific user
// @Tags SocialMedia
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {array} models.SocialMediaPost
// @Failure 404 {object} gin.H
func (c *SocialMediaController) GetAllSocialMediaPostByUserID(ctx *gin.Context) {
	// Get userID from URL
	userID := ctx.Param("user_id")

	// Get posts details by userID
	post, err := c.socialmediaService.GetAllSocialMediaPostByUserID(userID)
	if err != nil {
		c.logger.WithError(err).Error("Failed to get social media posts")
		ctx.JSON(http.StatusNotFound, gin.H{"status": "Post Not Found"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"post": post})
}

// GetSocialMediaPostByPostAndUserID retrieves a post by PostID and UserID
// @Summary Retrieve a social media post by PostID and UserID
// @Description Fetches a single social media post by its PostID and the UserID of the creator
// @Tags SocialMedia
// @Accept json
// @Produce json
// @Param post_id path string true "Post ID"
// @Param user_id query string true "User ID"
// @Success 200 {object} models.SocialMediaPost
// @Failure 404 {object} gin.H
func (c *SocialMediaController) GetSocialMediaPostByPostAndUserID(ctx *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID := ctx.GetString("user_id") // current user login into the system

	// Get task ID from URL
	postID := ctx.Param("post_id")

	// Get task
	post, err := c.socialmediaService.GetSocialMediaPostByPostAndUserID(postID, userID)
	if err != nil {
		c.logger.WithError(err).Error("Failed to get social media post")
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"post": post})
}

// CreateSocialMediaPost creates a new social media post
// @Summary Create a new social media post
// @Description Creates a new social media post for the authenticated user
// @Tags SocialMedia
// @Accept json
// @Produce json
// @Param post body models.SocialMediaPost true "Social Media Post"
// @Success 201 {object} models.SocialMediaPost
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
func (c *SocialMediaController) CreateSocialMediaPost(ctx *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID := ctx.GetString("user_id")

	// Parse request body
	var post models.SocialMediaPost
	if err := ctx.ShouldBindJSON(&post); err != nil {
		c.logger.WithError(err).Error("Failed to parse request body")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create Social Media Post
	createdSocialMediaPost, err := c.socialmediaService.CreateSocialMediaPost(&post, userID)
	if err != nil {
		c.logger.WithError(err).Error("Failed to create post")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"post": createdSocialMediaPost})
}

// UpdateSocialMediaPost updates an existing social media post
// @Summary Update a social media post
// @Description Updates an existing social media post for the authenticated user
// @Tags SocialMedia
// @Accept json
// @Produce json
// @Param post_id path string true "Post ID"
// @Param post body models.SocialMediaPost true "Updated Social Media Post"
// @Success 200 {object} models.SocialMediaPost
// @Failure 400 {object} gin.H
// @Failure 404 {object} gin.H
func (c *SocialMediaController) UpdateSocialMediaPost(ctx *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID := ctx.GetString("user_id")

	// Get task ID from URL
	postID := ctx.Param("post_id")

	// Parse request body
	var post models.SocialMediaPost
	if err := ctx.ShouldBindJSON(&post); err != nil {
		c.logger.WithError(err).Error("Failed to parse request body")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update task
	updatedSocialMediaPost, err := c.socialmediaService.UpdateSocialMediaPost(postID, &post, userID)
	if err != nil {
		c.logger.WithError(err).Error("Failed to update social media post")
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"post": updatedSocialMediaPost})
}

// DeleteSocialMediaPost deletes a social media post
// @Summary Delete a social media post
// @Description Deletes a social media post by its ID
// @Tags SocialMedia
// @Accept json
// @Produce json
// @Param post_id path string true "Post ID"
// @Success 200 {object} gin.H
// @Failure 404 {object} gin.H
func (c *SocialMediaController) DeleteSocialMediaPost(ctx *gin.Context) {
	// Get task ID from URL
	postID := ctx.Param("post_id")

	// Delete post
	err := c.socialmediaService.DeleteSocialMediaPost(postID)
	if err != nil {
		c.logger.WithError(err).Error("Failed to delete social media post")
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}
