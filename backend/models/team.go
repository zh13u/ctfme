package models

import "gorm.io/gorm"

type Team struct {
    gorm.Model
    Name       string `gorm:"unique;not null"`
    InviteCode string `gorm:"unique;not null"`
    Users      []User
}
