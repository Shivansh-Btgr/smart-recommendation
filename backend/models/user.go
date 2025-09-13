package models

type User struct {
	ID                uint   `gorm:"primaryKey"`
	Email             string `gorm:"uniqueIndex;not null"`
	PasswordHash      string `gorm:"not null"`
	IsProfileComplete bool   `gorm:"default:false"`
	Profile           Profile
}

type Profile struct {
	ID               uint `gorm:"primaryKey"`
	UserID           uint `gorm:"uniqueIndex"`
	Name             string
	Phone            string
	Education        string
	CGPA             float32
	Skills           []string `gorm:"type:jsonb;serializer:json"`
	Experience       int
	SocialLinks      []string `gorm:"type:jsonb;serializer:json"`
	Location         string
	Interest         string
	ResumeLink       string
	PreferredJobType string
	Availability     string
	Languages        []string `gorm:"type:jsonb;serializer:json"`
}
