package services

import (
	"errors"
	"fmt"
	"strings"

	"go-azure/models"
	"go-azure/utils"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// SocialMediaService handles task operations
type SocialMediaService struct {
	db     *gorm.DB
	logger *logrus.Logger
}

// NewSocialMediaService creates a new SocialMediaService
func NewSocialMediaService() *SocialMediaService {
	return &SocialMediaService{
		db:     utils.GetDB(),
		logger: utils.GetLogger(),
	}
}

// GetAllSocialMediaPosts returns all posts with pagination and sorting
func (s *SocialMediaService) GetAllSocialMediaPosts(page int, pageSize int, sortBy, sortOrder string) (map[string]interface{}, error) {
	var posts []models.SocialMediaPost
	var totalCount int64

	// Defaults
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if sortBy == "" {
		sortBy = "created_at"
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}

	offset := (page - 1) * pageSize

	// Count total
	if err := s.db.Raw("SELECT COUNT(post_id) FROM social_media_posts").Scan(&totalCount).Error; err != nil {
		s.logger.WithError(err).Error("Count failed")
		return nil, err
	}

	// Fetch paginated posts
	if err := s.db.
		Select("post_id, post_text, post_image, user_id, created_at").
		Order(sortBy + " " + sortOrder).
		Limit(pageSize).
		Offset(offset).
		Find(&posts).Error; err != nil {
		s.logger.WithError(err).Error("Query failed")
		return nil, err
	}

	totalPages := (totalCount + int64(pageSize) - 1) / int64(pageSize)

	return map[string]any{
		"posts":        posts,
		"total_count":  totalCount,
		"current_page": page,
		"total_pages":  totalPages,
	}, nil
}

// New query with optimized query tested on 500 thousand records
//GUIDE FOR QUERY OPTIMIZATION ON mySQL table
// --Add a FULLTEXT index on post_text:

// ALTER TABLE social_media_posts ADD FULLTEXT(post_text);

// --Ensure created_at is indexed if sorting frequently:
// CREATE INDEX idx_created_at ON social_media_posts(created_at DESC);
// -- Ensure a FULLTEXT index exists on `post_text`
// SELECT *
// FROM social_media_posts
// WHERE MATCH(post_text) AGAINST ('*test*' IN NATURAL LANGUAGE MODE)
// ORDER BY created_at DESC
// LIMIT 10;
// OFFSET 0;
func (s *SocialMediaService) QuerySocialMediaPost(page, limit int, colName, searchText, sortBy, sortOrder string) (map[string]any, error) {
	var posts []models.SocialMediaPost
	var totalCount int64
	var filteredCount int64

	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	// Validate sort order
	validSortOrder := "asc"
	if strings.ToLower(sortOrder) == "desc" {
		validSortOrder = "desc"
	}

	// Restrict sortable columns
	validSortBy := "created_at"
	if sortBy != "" {
		switch sortBy {
		case "created_at", "updated_at", "likes", "post_text":
			validSortBy = sortBy
		}
	}

	// Base query
	query := s.db.Model(&models.SocialMediaPost{})

	// Total count before filtering
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Apply full-text search filter if column name has value - Wilcard search supported
	if colName == "post_text" && searchText != "" {
		wildcardQuery := "*" + searchText + "*" // Apply wildcard full-text search (BOOLEAN MODE)
		query = query.Where("MATCH("+colName+") AGAINST(? IN BOOLEAN MODE)", wildcardQuery)
	} else {
		query = query.Where(fmt.Sprintf("%s LIKE ?", colName), "%"+searchText+"%")
	}

	// Apply full-text search filter if post_text is searched not exact match
	// if colName == "post_image" && searchText != "" {
	// 	query = query.Where(fmt.Sprintf("%s LIKE ?", colName), "%"+searchText+"%")
	// }

	// Count after filtering
	if err := query.Count(&filteredCount).Error; err != nil {
		return nil, err
	}

	// Apply sort, limit, and offset
	if err := query.Order(fmt.Sprintf("%s %s", validSortBy, validSortOrder)).
		Limit(limit).
		Offset(offset).
		Find(&posts).Error; err != nil {
		return nil, err
	}

	totalPages := (filteredCount + int64(limit) - 1) / int64(limit)

	return map[string]any{
		"current_page":   page,
		"filtered_count": filteredCount,
		"posts":          posts,
		"total_count":    totalCount,
		"total_pages":    totalPages,
	}, nil
}

// Old query with very slow performance tested on 500 thousand records
func (s *SocialMediaService) QuerySocialMediaPost2(page, limit int, colName, searchText, sortBy, sortOrder string) (map[string]interface{}, error) {
	var posts []models.SocialMediaPost
	var totalCount, filteredCount int64

	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	if sortBy == "" {
		sortBy = "created_at"
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "asc"
	}
	offset := (page - 1) * limit

	// Count total
	if err := s.db.Raw("SELECT COUNT(post_id) FROM social_media_posts").Scan(&totalCount).Error; err != nil {
		s.logger.WithError(err).Error("Count failed")
		return nil, err
	}

	query := s.db.Model(&models.SocialMediaPost{})

	// Apply filter if applicable
	if colName != "" && searchText != "" {
		query = query.Where(colName+" LIKE ?", "%"+searchText+"%")
	}

	// Count filtered
	if err := query.Count(&filteredCount).Error; err != nil {
		s.logger.WithError(err).Error("Failed to count filtered")
		return nil, err
	}

	// Fetch filtered and paginated posts
	if err := query.
		// Select("post_id, user_id, post_text, post_image, created_at").
		// Where(colName+" LIKE ?", "%"+searchText+"%").
		Order(sortBy + " " + sortOrder).
		Limit(limit).
		Offset(offset).
		Find(&posts).Error; err != nil {
		s.logger.WithError(err).Error("Failed to query posts")
		return nil, err
	}

	totalPages := (filteredCount + int64(limit) - 1) / int64(limit)

	return map[string]interface{}{
		"posts":          posts,
		"total_count":    totalCount,
		"filtered_count": filteredCount,
		"current_page":   page,
		"total_pages":    totalPages,
	}, nil
}

// GetSocialMediaPostByID returns a post by PostID
func (s *SocialMediaService) GetSocialMediaPostByPostID(PostID string) (*models.SocialMediaPost, error) {
	var post models.SocialMediaPost

	result := s.db.Where("post_id = ?", PostID).Find(&post)
	if result.Error != nil {
		s.logger.WithError(result.Error).Error("Failed to get social media post")
		return nil, errors.New("post not found")
	}

	return &post, nil
}

// GetAllSocialMediaPost returns all post for a user
func (s *SocialMediaService) GetAllSocialMediaPostByUserID(userID string) ([]*models.SocialMediaPost, error) {
	var posts []*models.SocialMediaPost

	result := s.db.Where("user_id = ?", userID).Find(&posts)
	if result.Error != nil {
		s.logger.WithError(result.Error).Error("Failed to get social media posts")
		return nil, errors.New("failed to get social media posts")
	}

	return posts, nil
}

// GetSocialMediaPostByID returns a post by PostID and UserID
func (s *SocialMediaService) GetSocialMediaPostByPostAndUserID(PostID string, userID string) (*models.SocialMediaPost, error) {
	var post models.SocialMediaPost

	result := s.db.Where("post_id = ? AND user_id = ?", PostID, userID).First(&post)
	if result.Error != nil {
		s.logger.WithError(result.Error).Error("Failed to get social media post")
		return nil, errors.New("post not found")
	}

	return &post, nil
}

// CreateSocialMediaPost creates a new post
func (s *SocialMediaService) CreateSocialMediaPost(post *models.SocialMediaPost, userID string) (*models.SocialMediaPost, error) {
	// Set task ID and user ID
	post.PostID = uuid.New().String()
	post.UserID = userID

	// Create task in database
	result := s.db.Create(post)
	if result.Error != nil {
		s.logger.WithError(result.Error).Error("Failed to create task")
		return nil, errors.New("failed to create task")
	}

	s.logger.WithFields(logrus.Fields{
		"post_id": post.PostID,
		"user_id": userID,
	}).Info("Post created")

	return post, nil
}

// UpdateSocialMediaPost updates an existing post
func (s *SocialMediaService) UpdateSocialMediaPost(postID string, updatedSocialMediaPost *models.SocialMediaPost, userID string) (*models.SocialMediaPost, error) {
	// Get existing task
	var existingSocialMediaPost models.SocialMediaPost
	result := s.db.Where("post_id = ? AND user_id = ?", postID, userID).First(&existingSocialMediaPost)
	if result.Error != nil {
		s.logger.WithError(result.Error).Error("Failed to get social media post for update")
		return nil, errors.New("post not found")
	}

	// Update post fields
	existingSocialMediaPost.PostText = updatedSocialMediaPost.PostText
	existingSocialMediaPost.PostImage = updatedSocialMediaPost.PostImage

	// Save changes to database
	result = s.db.Save(&existingSocialMediaPost)
	if result.Error != nil {
		s.logger.WithError(result.Error).Error("Failed to update social media post")
		return nil, errors.New("failed to update social media post")
	}

	s.logger.WithFields(logrus.Fields{
		"post_id": postID,
		"user_id": userID,
	}).Info("Task updated")

	return &existingSocialMediaPost, nil
}

// DeleteSocialMediaPost deletes a task
func (s *SocialMediaService) DeleteSocialMediaPost(postID string) error {
	// Check if task exists and belongs to user
	var post models.SocialMediaPost
	result := s.db.Where("post_id = ?", postID).First(&post)
	if result.Error != nil {
		s.logger.WithError(result.Error).Error("Failed to get social media post for deletion")
		return errors.New("post not found")
	}

	// Delete post
	result = s.db.Delete(&post)
	if result.Error != nil {
		s.logger.WithError(result.Error).Error("Failed to delete social media post")
		return errors.New("failed to delete social media post")
	}

	s.logger.WithFields(logrus.Fields{
		"post_id": postID,
		// "user_id": userID,
	}).Info("Post deleted")

	return nil
}
