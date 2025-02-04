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
	Username string    `json:"username" gorm:"unique;not null"`
	Email    string    `json:"email" gorm:"unique;not null"`
	Password string    `json:"password" gorm:"not null"`
	Roles    []Role    `gorm:"many2many:user_roles"`
}

type AuthResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Token string `json:"token"`
}

