package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"sas-pro/internal/models"
	"sas-pro/internal/repositories"
)

func CreateProduct(c *gin.Context) {
	var input models.ProductCreateRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	repo := repositories.NewProductRepository()
	product, err := repo.CreateProductWithVariants(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, product)
}

func GetProducts(c *gin.Context) {
	repo := repositories.NewProductRepository()
	products, err := repo.GetAllProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}