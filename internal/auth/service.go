package auth

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	_ "time"
)

type AuthService struct {
	db *gorm.DB
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{db: db}
}

func (s *AuthService) GenerateOTP() (string, error) {
	buffer := make([]byte, 6)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}
	// print the OTP to the console
	fmt.Println(base32.StdEncoding.EncodeToString(buffer)[:6])
	return base32.StdEncoding.EncodeToString(buffer)[:6], nil
}

func (s *AuthService) SendOTP(phoneNumber, otp string) error {
	// Implement SMS sending logic here
	return nil
}

func (s *AuthService) VerifyOTP(phoneNumber, otp string) bool {
	var user User
	if err := s.db.Where("phone_number = ?", phoneNumber).First(&user).Error; err != nil {
		return false
	}
	return user.OTPSecret == otp
}

func (s *AuthService) SetPIN(phoneNumber, pin string) error {
	hashedPIN, err := bcrypt.GenerateFromPassword([]byte(pin), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.db.Model(&User{}).Where("phone_number = ?", phoneNumber).
		Update("pin", string(hashedPIN)).Error
}

func (s *AuthService) VerifyPIN(phoneNumber, pin string) bool {
	var user User
	if err := s.db.Where("phone_number = ?", phoneNumber).First(&user).Error; err != nil {
		return false
	}
	return bcrypt.CompareHashAndPassword([]byte(user.PIN), []byte(pin)) == nil
}
