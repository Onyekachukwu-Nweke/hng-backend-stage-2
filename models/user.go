package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	UserID    string `gorm:"primaryKey" json:"userId"`
	FirstName string `gorm:"not null" json:"firstName" validate:"required"`
	LastName  string `gorm:"not null" json:"lastName" validate:"required"`
	Email     string `gorm:"unique;not null" json:"email" validate:"required,email"`
	Password  string `gorm:"not null" json:"password" validate:"required"`
	Phone     string `json:"phone"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.UserID = uuid.New().String()
	return
}
