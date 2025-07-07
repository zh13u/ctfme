package models

import "gorm.io/gorm"

type Setup struct {
	gorm.Model
    CTFMode             string `gorm:"default:'user'"` 
    DynamicScoreEnabled bool   `gorm:"default:false"`
    DynamicScoreDecay   int    `gorm:"default:10"`
    DynamicScoreMin     int    `gorm:"default:50"`
}
