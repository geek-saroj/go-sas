package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name        string     `gorm:"size:255;not null"`
	Description string     `gorm:"size:1000"`
	Variants    []Variant  `gorm:"foreignKey:ProductID"`
}

type Variant struct {
	gorm.Model
	ProductID     uint           `gorm:"not null"`
	Name          string         `gorm:"size:255;not null"`
	SerialNumbers []SerialNumber `gorm:"foreignKey:VariantID"`
}

type SerialNumber struct {
	gorm.Model
	VariantID  uint   `gorm:"not null"`
	Number     string `gorm:"size:255;uniqueIndex;not null"`
	IsUsed     bool   `gorm:"default:false"`
}

type ProductCreateRequest struct {
	Name        string             `json:"name" binding:"required"`
	Description string             `json:"description"`
	Variants    []VariantCreateDTO `json:"variants" binding:"required"`
}

type VariantCreateDTO struct {
	Name    string   `json:"name" binding:"required"`
	Serials []string `json:"serials" binding:"required"`
}