package handlers

import (
	"net/http"
	"strconv"
	"time"

	"backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)


type InternshipInput struct {
	Title        string    `json:"title" binding:"required"`
	Organization string    `json:"organization" binding:"required"`
	Location     string    `json:"location" binding:"required"`
	StipendINR   int       `json:"stipend_inr" binding:"required"`
	Duration     int       `json:"duration"`
	SkillsReq    []string  `json:"skillsreq"`
	LangsReq     []string  `json:"langsreq"`
	Active       *bool     `json:"active"`
	PostedAt     time.Time `json:"posted_at"`
	Deadline     time.Time `json:"deadline" binding:"required"`
	ApplyURL     string    `json:"apply_url" binding:"required"`
	Description  string    `json:"description"`
	MinCGPA      float32   `json:"min_cgpa"`
	Experience   int       `json:"experience"`
}

func CreateInternship(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input InternshipInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		   internship := models.Internship{
			   Title:        input.Title,
			   Organization: input.Organization,
			   Location:     input.Location,
			   StipendINR:   input.StipendINR,
			   Duration:     input.Duration,
			   SkillsReq:    input.SkillsReq,
			   LangsReq:     input.LangsReq,
			   Active:       input.Active,
			   PostedAt:     time.Now(),
			   Deadline:     input.Deadline,
			   ApplyURL:     input.ApplyURL,
			   Description:  input.Description,
			   MinCGPA:      input.MinCGPA,
			   Experience:   input.Experience,
		   }
		if err := db.Create(&internship).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create internship", "details": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, internship)
	}
}

func EditInternship(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid internship ID"})
			return
		}
		var input InternshipInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var internship models.Internship
		if err := db.First(&internship, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Internship not found"})
			return
		}
		internship.Title = input.Title
		internship.Organization = input.Organization
		internship.Location = input.Location
		internship.StipendINR = input.StipendINR
		internship.Duration = input.Duration
		internship.SkillsReq = input.SkillsReq
		internship.LangsReq = input.LangsReq
		internship.Active = input.Active
		internship.Deadline = input.Deadline
		internship.ApplyURL = input.ApplyURL
		internship.Description = input.Description
		internship.MinCGPA = input.MinCGPA
		internship.Experience = input.Experience
		if err := db.Save(&internship).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update internship", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, internship)
	}
}

func GetAllInternships(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var internships []models.Internship
		query := db.Model(&models.Internship{})

		// Filtering
		if loc := c.Query("location"); loc != "" {
			query = query.Where("location ILIKE ?", "%"+loc+"%")
		}
		if minStipend := c.Query("min_stipend"); minStipend != "" {
			if val, err := strconv.Atoi(minStipend); err == nil {
				query = query.Where("stipend_inr >= ?", val)
			}
		}
		if minDuration := c.Query("min_duration"); minDuration != "" {
			if val, err := strconv.Atoi(minDuration); err == nil {
				query = query.Where("duration >= ?", val)
			}
		}
		if active := c.Query("active"); active != "" {
			if active == "true" {
				query = query.Where("active = ?", true)
			} else if active == "false" {
				query = query.Where("active = ?", false)
			}
		}

		// Sorting
		sort := c.DefaultQuery("sort", "")
		if sort == "stipend" {
			query = query.Order("stipend_inr DESC")
		} else if sort == "deadline" {
			query = query.Order("deadline ASC")
		} else {
			query = query.Order("posted_at DESC")
		}

		if err := query.Find(&internships).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch internships"})
			return
		}
		c.JSON(http.StatusOK, internships)
	}
}

func GetInternshipByID(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid internship ID"})
			return
		}
		var internship models.Internship
		if err := db.First(&internship, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Internship not found"})
			return
		}
		c.JSON(http.StatusOK, internship)
	}
}

func GetActiveInternships(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var internships []models.Internship
		if err := db.Where("active = ?", true).Find(&internships).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch active internships"})
			return
		}
		c.JSON(http.StatusOK, internships)
	}
}
