package repositories

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ServiceRepository struct {
	db *pgxpool.Pool
}

func NewServiceRepository(db *pgxpool.Pool) *ServiceRepository {
	return &ServiceRepository{db: db}
}

func (r *ServiceRepository) GetURL(serviceType string) (string, error) {
	var url string
	err := r.db.QueryRow(
		context.Background(),
		"SELECT url FROM service_types WHERE name = $1", 
		serviceType,
	).Scan(&url)
	
	if err != nil {
		return "", fmt.Errorf("service type not found: %v", serviceType)
	}
	return url, nil
}