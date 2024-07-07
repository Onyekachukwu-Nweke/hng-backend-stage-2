package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Organisation struct {
	OrgID       string `gorm:"primaryKey" json:"orgId"`
	Name        string `gorm:"not null" json:"name" validate:"required"`
	Description string `json:"description"`
}

func (o *Organisation) BeforeCreate(tx *gorm.DB) (err error) {
	o.OrgID = uuid.New().String()
	return
}
