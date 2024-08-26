package models

import (
	"gorm.io/gorm"
)

type Users struct {
	gorm.Model
	ID       int
	Name     string
	Email    string
	Password []byte
}

type Sessions struct {
	gorm.Model
	ID           string
	UserEmail    string
	UserIP       string
	RefreshToken string
	Attemp       int
}
