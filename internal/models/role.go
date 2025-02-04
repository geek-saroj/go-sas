package models

import (
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	Name        string         `json:"name" gorm:"unique;not null"`
	Permissions []Permission  `gorm:"many2many:role_permissions"`
	Users       []User        `gorm:"many2many:user_roles"`
}
