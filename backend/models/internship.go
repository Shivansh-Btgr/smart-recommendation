package models

import (
	"time"
)

type Internship struct {
	ID           uint      `gorm:"primaryKey"`
	Title        string    `gorm:"not null"`
	Organization string    `gorm:"not null"`
	Location     string    `gorm:"not null"`
	StipendINR   int       `gorm:"not null"`
	Duration     int       `gorm:"not null"`
	SkillsReq    []string  `gorm:"type:jsonb;serializer:json"`
	LangsReq     []string  `gorm:"type:jsonb;serializer:json"`
	Active       *bool     `gorm:"default:true"`
	PostedAt     time.Time `gorm:"autoCreateTime"`
	Deadline     time.Time
	ApplyURL     string  `gorm:"not null"`
	Description  string
	MinCGPA      float32 
	Experience   int
}
