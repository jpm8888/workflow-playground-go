package auth

import "gorm.io/gorm"

type User struct {
	gorm.Model
	PhoneNumber string `gorm:"uniqueIndex"`
	PIN         string
	OTPSecret   string
}
