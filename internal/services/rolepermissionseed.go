package services

import (
	"sas-pro/internal/models"

	"gorm.io/gorm"
)

func SeedRolesAndPermissions(db *gorm.DB) {
	// Seed permissions
	createProductPermission := models.Permission{
		Name:        "CreateProduct",
		Description: "Permission to create a product",
	}
	db.Create(&createProductPermission)

	// Seed roles
	AdminRole := models.Role{
		Name:        "Admin",
		Permissions: []models.Permission{createProductPermission},
	}
	db.Create(&AdminRole)
}
