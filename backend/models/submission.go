package models

import "gorm.io/gorm"

type Submission struct {
    gorm.Model
    UserID      uint
    TeamID      *uint
    ChallengeID uint
    Flag        string
    IsCorrect   bool
}
