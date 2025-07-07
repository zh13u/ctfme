package models

import "gorm.io/gorm"

type Challenge struct {
	gorm.Model

	Title       string `gorm:"not null"`
	Description string `gorm:"type:text"`
	Category    string `gorm:"not null"`
	Points      int    `gorm:"not null"`
	Flag        string `gorm:"not null"`
	FileURL     string
	Visible     bool   `gorm:"default:true"`
	SolvedBy    []User `gorm:"many2many:challenge_solvers;"`
}
