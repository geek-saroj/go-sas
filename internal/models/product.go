package models

import (
	"mime/multipart"

	"gorm.io/gorm"
)

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
    VariablesFile   *multipart.FileHeader `form:"variables_file"`  // for the file upload
    Direction       string                `form:"direction,omitempty"`
    VariableName    string                `form:"variable_name,omitempty"`
    VariableType    string                `form:"variable_type,omitempty"`
    DefaultValue    string                `form:"default_value"`
    UseDefaultValue string                `form:"use_default_value"`
    LinkSlug        string                `form:"link_slug,omitempty"`
}


type VariantCreateDTO struct {
	Name    string   `json:"name" binding:"required"`
	Serials []string `json:"serials" binding:"required"`
}