package routes

import (
	"sas-pro/internal/handlers"
	middlewares "sas-pro/internal/handlers/middleware"

	"github.com/gin-gonic/gin"
)

func Setup(router *gin.Engine) {
	auth := router.Group("/auth")
	{
		auth.POST("/register", handlers.Register)
		auth.POST("/login", handlers.Login)
		auth.POST("/permissions", handlers.CreatePermissions)        
		auth.POST("/roles", handlers.CreateRoleAndPermissions)                  
		// auth.POST("/roles/:role_id/permissions/:permission_id", handlers.AssignPermissionToRole)
		auth.POST("/users/:user_id/roles/:role_id", handlers.AssignRoleToUser)                  
	}


	api := router.Group("/api")
	
	{
		products := api.Group("/products")
		{
			products.POST("",middlewares.CheckPermission("createuser"), handlers.CreateProduct)
			// products.GET("",middlewares.CheckPermission("createuser"), handlers.GetProducts)
		}
	}
}