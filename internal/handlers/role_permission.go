package handlers

import (
	"net/http"
	"sas-pro/internal/models"
	"sas-pro/pkg/database"

	"github.com/gin-gonic/gin"
)

// Initialize the DB connection in the main app

// Create a new permission
func CreatePermission(c *gin.Context) {
	var permission models.Permission
	if err := c.ShouldBindJSON(&permission); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	// Check if permission already exists
	if err := database.DB.Where("name = ?", permission.Name).First(&models.Permission{}).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Permission already exists"})
		return
	}

	if err := database.DB.Create(&permission).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create permission"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": permission})
}

// Create a new role
func CreateRole(c *gin.Context) {
	var role models.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	// Check if role already exists
	if err := database.DB.Where("name = ?", role.Name).First(&models.Role{}).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role already exists"})
		return
	}

	if err := database.DB.Create(&role).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create role"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": role})
}

// Assign a permission to a role
func AssignPermissionToRole(c *gin.Context) {
	roleID := c.Param("role_id")
	permissionID := c.Param("permission_id")

	var role models.Role
	var permission models.Permission

	// Fetch the role and permission
	if err := database.DB.First(&role, roleID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
		return
	}

	if err := database.DB.First(&permission, permissionID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Permission not found"})
		return
	}

	// Assign permission to role
	if err := database.DB.Model(&role).Association("Permissions").Append(&permission); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign permission to role"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Permission assigned to role successfully"})
}

// Assign a role to a user
func AssignRoleToUser(c *gin.Context) {
	userID := c.Param("user_id")
	roleID := c.Param("role_id")

	var user models.User
	var role models.Role

	// Fetch the user and role
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if err := database.DB.First(&role, roleID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
		return
	}

	// Assign role to user
	if err := database.DB.Model(&user).Association("Roles").Append(&role); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign role to user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role assigned to user successfully"})
}
