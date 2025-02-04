package routes

import (
	"sas-pro/internal/handlers"
	middlewares "sas-pro/internal/handlers/middleware"

	"github.com/gin-gonic/gin"
)

func Setup(router *gin.Engine) {
	// Auth routes
	auth := router.Group("/auth")
	{
		auth.POST("/register", handlers.Register)
		auth.POST("/login", handlers.Login)
		auth.POST("/permissions", handlers.CreatePermissions)         // Create permission
		auth.POST("/roles", handlers.CreateRoleAndPermissions)                     // Create role
		// auth.POST("/roles/:role_id/permissions/:permission_id", handlers.AssignPermissionToRole) // Assign permission to role
		auth.POST("/users/:user_id/roles/:role_id", handlers.AssignRoleToUser)                   // Assign role to user
	}

	// Authenticated routes
	api := router.Group("/api")
	
	// api.Use(middleware.JWTAuth())
	 api.POST("/products", middlewares.CheckPermission("createuser"), handlers.CreateProduct)
	// {
	// 	products := api.Group("/products")
	// 	{
	// 		products.POST("", handlers.CreateProduct)
	// 		products.GET("", handlers.GetProducts)
	// 	}
	// }
}