package main

import (
	"backend/handlers"
	"backend/models"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	_ = godotenv.Load()

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbSSL := os.Getenv("DB_SSLMODE")
	if dbSSL == "" {
		dbSSL = "disable"
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", dbHost, dbUser, dbPass, dbName, dbPort, dbSSL)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to database: ", err)
	}

	db.AutoMigrate(&models.User{}, &models.Profile{}, &models.Internship{})

	r := gin.Default()

	r.POST("/register", handlers.Register(db))
	r.POST("/login", handlers.Login(db))

	r.GET("/internships", handlers.GetAllInternships(db))
	r.GET("/internships/:id", handlers.GetInternshipByID(db))
	r.GET("/internships/active", handlers.GetActiveInternships(db))

	protected := r.Group("/")
	protected.Use(handlers.AuthMiddleware(db))
	{
		protected.GET("/profile", handlers.GetProfile(db))
		protected.PUT("/profile", handlers.UpdateProfile(db))
		mlURL := os.Getenv("ML_SERVICE_URL")
		protected.GET("/recommendations", handlers.GetRecommendations(db, mlURL))
		protected.POST("/internships", handlers.CreateInternship(db))
		protected.PUT("/internships/:id", handlers.EditInternship(db))
	}

	r.Run(":8080")
}
