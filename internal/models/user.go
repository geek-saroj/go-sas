package models

import "gorm.io/gorm"

// type User struct {
// 	gorm.Model
// 	Email    string `gorm:"uniqueIndex;not null"`
// 	Password string `gorm:"not null"`
// 	RoleID   uint   `gorm:"not null "`
// }

type User struct {
	gorm.Model
	Username string `json:"username" gorm:"unique;not null" validate:"required,min=3,max=20"`
	Email    string `json:"email" gorm:"unique;not null" validate:"required,email"`
	Password string `json:"password" gorm:"not null" validate:"required,min=6,max=100"`
	Roles    []Role `gorm:"many2many:user_roles"`
}

type AuthResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Token string `json:"token"`
}
