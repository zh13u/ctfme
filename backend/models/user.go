package models

import "gorm.io/gorm"

type User struct {
	gorm.Model

	Username    string `gorm:"unique;not null"`
	Email       string `gorm:"unique;not null"`
	Password    string `gorm:"not null"`
	TeamID      *uint
	IsAdmin     bool `gorm:"default:false"`
	Team        *Team
	Submissions []Submission
}
