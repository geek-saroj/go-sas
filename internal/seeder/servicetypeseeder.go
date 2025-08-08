package seeder

import (
	"log"
	"sas-pro/internal/models"

	"gorm.io/gorm"
)

// ServiceTypeSeeder will seed the ServiceType data into the database
type ServiceTypeSeeder struct {
	DB *gorm.DB
}

// NewServiceTypeSeeder creates a new instance of the ServiceTypeSeeder
func NewServiceTypeSeeder(db *gorm.DB) *ServiceTypeSeeder {
	return &ServiceTypeSeeder{DB: db}
}

// Seed will insert ServiceType records into the database
func (s *ServiceTypeSeeder) Seed() {
	serviceTypes := []models.ServiceType{
		{Name: "auth", Description: "Authentication service", Url: "http://localhost:5000/users"},
		{Name: "data", Description: "Payment gateway service", Url: "http://localhost:500/data"},
		{Name: "user", Description: "User management service", Url: "http://user-service/api"},
		{Name: "order", Description: "Order management service", Url: "http://order-service/api"},
	}

	for _, serviceType := range serviceTypes {
		// Check if the ServiceType already exists based on the Name (uniqueIndex)
		if err := s.DB.Where("name = ?", serviceType.Name).First(&models.ServiceType{}).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				// Create new record if it doesn't exist
				if err := s.DB.Create(&serviceType).Error; err != nil {
					log.Printf("Error seeding ServiceType '%s': %v", serviceType.Name, err)
				} else {
					log.Printf("Seeded ServiceType '%s'", serviceType.Name)
				}
			} else {
				log.Printf("Error checking for ServiceType '%s': %v", serviceType.Name, err)
			}
		}
	}
}
