package routes

import (
	"sas-pro/internal/handlers"
	"sas-pro/internal/handlers/middleware"
	"github.com/gin-gonic/gin"
)

func Setup(router *gin.Engine) {
	// Auth routes
	auth := router.Group("/auth")
	{
		auth.POST("/register", handlers.Register)
		auth.POST("/login", handlers.Login)
	}

	// Authenticated routes
	api := router.Group("/api")
	api.Use(middleware.JWTAuth())
	{
		products := api.Group("/products")
		{
			products.POST("", handlers.CreateProduct)
			products.GET("", handlers.GetProducts)
		}
	}
}