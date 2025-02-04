package repositories

import (
	"sas-pro/internal/models"
	"sas-pro/pkg/database"

	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository() *ProductRepository {
	return &ProductRepository{db: database.DB}
}

func (r *ProductRepository) CreateProductWithVariants(input models.ProductCreateRequest) (*models.Product, error) {
	product := models.Product{
		Name:        input.Name,
		Description: input.Description,
	}

	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&product).Error; err != nil {
			return err
		}

		for _, variantDTO := range input.Variants {
			variant := models.Variant{
				ProductID: product.ID,
				Name:      variantDTO.Name,
			}

			if err := tx.Create(&variant).Error; err != nil {
				return err
			}

			var serials []models.SerialNumber
			for _, s := range variantDTO.Serials {
				serials = append(serials, models.SerialNumber{
					VariantID: variant.ID,
					Number:    s,
				})
			}

			if err := tx.Create(&serials).Error; err != nil {
				return err
			}
		}

		return nil
	})

	return &product, err
}

func (r *ProductRepository) GetAllProducts() ([]models.Product, error) {
	var products []models.Product
	err := r.db.Preload("Variants.SerialNumbers").Find(&products).Error
	return products, err
}