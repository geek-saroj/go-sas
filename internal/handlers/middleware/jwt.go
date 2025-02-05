package middlewares

import (
	"net/http"
	"sas-pro/config"
	"sas-pro/internal/models"
	"sas-pro/pkg/database"
	"sas-pro/pkg/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func CheckPermission(permissionName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		cfg := config.MustLoad()
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Bearer token required"})
			return
		}

		userID, err := utils.ParseJWT(tokenString, cfg.JWTSecret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}
		var user models.User

		db := database.DB
		if err := db.Preload("Roles.Permissions").First(&user, userID).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		hasPermission := false
		for _, role := range user.Roles {
			for _, perm := range role.Permissions {
				if perm.Name == permissionName {
					hasPermission = true
					break
				}
			}
			if hasPermission {
				break
			}
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
			c.Abort()
			return
		}

		c.Next()
	}
}
