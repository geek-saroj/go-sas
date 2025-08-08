package models

import (
	"gorm.io/gorm"
)

type ServiceType struct {
	gorm.Model
	Name        string `gorm:"uniqueIndex;not null"`
	Description string
	Url         string
	
}