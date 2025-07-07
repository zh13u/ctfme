package models

import "gorm.io/gorm"

type Submission struct {
	gorm.Model
	UserID       uint
	TeamID       *uint
	ChallengeID  uint
	Flag         string
	IsCorrect    bool
	PointsEarned int `gorm:"default:0"` // Points earned at the time of submission
}
