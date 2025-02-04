// // internal/handlers/middleware/jwt.go
// package middleware

// import (
// 	"net/http"
// 	"strings"

// 	"github.com/gin-gonic/gin"
// 	"sas-pro/pkg/utils"
// 	"sas-pro/config"
// )

// func JWTAuth() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		cfg := config.MustLoad()

// 		authHeader := c.GetHeader("Authorization")
// 		if authHeader == "" {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
// 			return
// 		}

// 		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
// 		if tokenString == authHeader {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Bearer token required"})
// 			return
// 		}

// 		userID, err := utils.ParseJWT(tokenString, cfg.JWTSecret)
// 		if err != nil {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
// 			return
// 		}

// 		// Set user ID in context for downstream handlers
// 		c.Set("userID", userID)
// 		c.Next()
// 	}
// }

package middlewares

import (
	"log"
	"net/http"
	"sas-pro/config"
	"sas-pro/internal/models"
	"sas-pro/pkg/database"
	"sas-pro/pkg/utils"
	"strings"

	"github.com/gin-gonic/gin"
	// "gorm.io/gorm"
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

		// Set user ID in context for downstream handlers
	
	  // Assume user ID is set after authentication
		var user models.User

		// Fetch user with roles and permissions
		log.Println("User ID from context:", userID)
		// userID := c.GetInt("user_id") 
		db := database.DB
		if err := db.Preload("Roles.Permissions").First(&user, userID).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		// Check if user has the required permission
		hasPermission := false
		for _, role := range user.Roles {
			for _, perm := range role.Permissions {
				log.Println("Role:", role.Name, "Permission:", perm.Name)
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
