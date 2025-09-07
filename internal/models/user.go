package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UserName     string `gorm:"size:100;unique"`
	PasswordHash string `gorm:"size:255"`
	Email        string `gorm:"size:255;unique"`
	Role         string `gorm:"size:50"`
}
