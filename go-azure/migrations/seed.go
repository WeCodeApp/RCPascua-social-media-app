package migrations

import (
	"math/rand"
	"strconv"
	"time"

	"go-azure/models"

	"github.com/bxcodec/faker/v3"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Seed populates the database with fake data
func Seed(db *gorm.DB) error {
	logrus.Info("Seeding database")

	// Seed users
	users, err := seedUsers(db, 10)
	if err != nil {
		return err
	}

	// Seed tasks for each user
	for _, user := range users {
		if err := seedTasks(db, user, 5); err != nil {
			return err
		}
	}

	// Seed social media posts for each user
	for _, user := range users {
		if err := seedSocialMediaPost(db, user, 5); err != nil {
			return err
		}
	}

	// Seed social media comments and likes for each post
	for _, user := range users { // Fetch all social media posts for each user
		var posts []models.SocialMediaPost
		if err := db.Where("user_id = ?", user.ID).Find(&posts).Error; err != nil {
			logrus.WithError(err).Error("Failed to fetch posts for user")
			return err
		}

		for _, post := range posts {
			// Seed social media comments for each post
			if err := seedSocialMediaComments(db, post, 5); err != nil {
				return err
			}
			// Seed social media likes for each post
			if err := seedSocialMediaLikes(db, post, 5); err != nil {
				return err
			}
		}
	}

	logrus.Info("Database seeding completed successfully")
	return nil
}

// seedUsers creates fake users
func seedUsers(db *gorm.DB, count int) ([]models.User, error) {
	var users []models.User

	for i := 0; i < count; i++ {
		user := models.User{
			ID:        uuid.New().String(),
			Email:     faker.Email(),
			Name:      faker.Name(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := db.Create(&user).Error; err != nil {
			logrus.WithError(err).Error("Failed to seed user")
			return nil, err
		}

		users = append(users, user)
	}

	logrus.WithField("count", count).Info("Users seeded successfully")
	return users, nil
}

// seedTasks creates fake tasks for a user
func seedTasks(db *gorm.DB, user models.User, count int) error {
	for i := 0; i < count; i++ {
		completed := false
		if i%2 == 0 {
			completed = true
		}

		task := models.Task{
			ID:          uuid.New().String(),
			Title:       faker.Sentence(),
			Description: faker.Paragraph(),
			Completed:   completed,
			UserID:      user.ID,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		if err := db.Create(&task).Error; err != nil {
			logrus.WithError(err).Error("Failed to seed task")
			return err
		}
	}

	logrus.WithFields(logrus.Fields{
		"user_id": user.ID,
		"count":   count,
	}).Info("Tasks seeded successfully")
	return nil
}

func RandomInt(min, max int) int {
	return min + rand.Intn(max-min+1)
}

func RandomBool() bool {
	return rand.Intn(2) == 0
}

func seedSocialMediaPost(db *gorm.DB, user models.User, count int) error {
	// Seed social media posts for the user
	for i := 0; i < count; i++ {
		imageLink := "https://picsum.photos/id/" + strconv.Itoa(RandomInt(1, 1000)) + "/" + strconv.Itoa(RandomInt(1, 1000)) + "/" + strconv.Itoa(RandomInt(1, 1000)) // Replace with a valid image URL
		post := models.SocialMediaPost{
			PostID:    uuid.New().String(),
			UserID:    user.ID,
			PostText:  faker.Sentence(),
			PostImage: imageLink, //faker.URL(),
			Likes:     RandomInt(1, 1000),
			IsLiked:   RandomInt(0, 1) == 1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := db.Create(&post).Error; err != nil {
			logrus.WithError(err).Error("Failed to seed social media post")
			return err
		}
	}
	logrus.WithField("count", count).Info("Post seeded successfully")
	return nil
}

func seedSocialMediaComments(db *gorm.DB, post models.SocialMediaPost, count int) error {
	// Seed social media comments for the post
	for i := 0; i < count; i++ {
		comment := models.SocialMediaComments{
			CommentID:   uuid.New().String(),
			PostID:      post.PostID,
			UserID:      post.UserID,
			CommentText: faker.Sentence(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		if err := db.Create(&comment).Error; err != nil {
			logrus.WithError(err).Error("Failed to seed social media comment")
			return err
		}
	}
	logrus.WithField("count", count).Info("Comments seeded successfully")
	return nil
}

func seedSocialMediaLikes(db *gorm.DB, post models.SocialMediaPost, count int) error {
	// Seed social media likes for the post
	for i := 0; i < count; i++ {
		like := models.SocialMediaLikes{
			LikeID:    uuid.New().String(),
			PostID:    post.PostID,
			UserID:    post.UserID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := db.Create(&like).Error; err != nil {
			logrus.WithError(err).Error("Failed to seed social media like")
			return err
		}
	}
	logrus.WithField("count", count).Info("Likes seeded successfully")
	return nil
}
