package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"backend/models"
)

type ProfileInput struct {
	Name             string   `json:"name"`
	Phone            string   `json:"phone"`
	Education        string   `json:"education"`
	CGPA             float32  `json:"cgpa"`
	Skills           []string `json:"skills"`
	Experience       int      `json:"experience"`
	SocialLinks      []string `json:"social_links"`
	Location         string   `json:"location"`
	Interest         string   `json:"interest"`
	ResumeLink       string   `json:"resume_link"`
	PreferredJobType string   `json:"preferred_job_type"`
	Availability     string   `json:"availability"`
	Languages        []string `json:"languages"`
}

func GetProfile(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet("user").(models.User)
		var profile models.Profile
		if err := db.Where("user_id = ?", user.ID).First(&profile).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
			return
		}
		c.JSON(http.StatusOK, profile)
	}
}

func UpdateProfile(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet("user").(models.User)
		var input ProfileInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var profile models.Profile
		if err := db.Where("user_id = ?", user.ID).First(&profile).Error; err == nil {
			// Update existing profile
			if err := db.Model(&profile).Updates(models.Profile{
				Name:             input.Name,
				Phone:            input.Phone,
				Education:        input.Education,
				CGPA:             input.CGPA,
				Skills:           input.Skills,
				Experience:       input.Experience,
				SocialLinks:      input.SocialLinks,
				Location:         input.Location,
				Interest:         input.Interest,
				ResumeLink:       input.ResumeLink,
				PreferredJobType: input.PreferredJobType,
				Availability:     input.Availability,
				Languages:        input.Languages,
			}).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile", "details": err.Error()})
				return
			}
			db.Where("user_id = ?", user.ID).First(&profile)
		} else {
			profile = models.Profile{
				UserID:           user.ID,
				Name:             input.Name,
				Phone:            input.Phone,
				Education:        input.Education,
				CGPA:             input.CGPA,
				Skills:           input.Skills,
				Experience:       input.Experience,
				SocialLinks:      input.SocialLinks,
				Location:         input.Location,
				Interest:         input.Interest,
				ResumeLink:       input.ResumeLink,
				PreferredJobType: input.PreferredJobType,
				Availability:     input.Availability,
				Languages:        input.Languages,
			}
			if err := db.Create(&profile).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create profile", "details": err.Error()})
				return
			}
		}
		if err := db.Model(&models.User{}).Where("id = ?", user.ID).Update("is_profile_complete", true).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user profile status", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Profile updated", "profile": profile})
	}
}
