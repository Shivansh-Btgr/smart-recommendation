package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"backend/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MLRecommendationRequest struct {
	Profile models.Profile `json:"profile"`
}

type MLRecommendationResponse struct {
	InternshipIDs []int `json:"internship_ids"`
}

func GetRecommendations(db *gorm.DB, mlURL string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet("user").(models.User)
		if !user.IsProfileComplete {
			c.JSON(http.StatusForbidden, gin.H{"error": "Complete your profile to get recommendations"})
			return
		}
		var profile models.Profile
		if err := db.Where("user_id = ?", user.ID).First(&profile).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
			return
		}
		reqBody, err := json.Marshal(MLRecommendationRequest{Profile: profile})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encode profile"})
			return
		}
		resp, err := http.Post(mlURL, "application/json", bytes.NewBuffer(reqBody))
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to contact ML service", "details": err.Error()})
			return
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		if resp.StatusCode != 200 {
			c.JSON(http.StatusBadGateway, gin.H{"error": "ML service error", "details": string(body)})
			return
		}
		var mlResp MLRecommendationResponse
		if err := json.Unmarshal(body, &mlResp); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse ML response"})
			return
		}
		var internships []models.Internship
		if len(mlResp.InternshipIDs) > 0 {
			if err := db.Where("id IN ?", mlResp.InternshipIDs).Find(&internships).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch recommended internships"})
				return
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"recommendation_ids": mlResp.InternshipIDs,
			"recommendations":    internships,
		})
	}
}
