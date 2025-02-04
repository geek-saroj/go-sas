package repositories

import (
	"sas-pro/internal/models"
	"sas-pro/pkg/database"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository() *UserRepository {
	return &UserRepository{db: database.DB}
}

func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}
func (r *UserRepository) HasPermission(userID uint, permission string) (bool, error) {
    var count int64
    err := r.db.Model(&models.User{}).
        Joins("JOIN user_roles ON user_roles.user_id = users.id").
        Joins("JOIN roles ON roles.id = user_roles.role_id").
        Joins("JOIN role_permissions ON role_permissions.role_id = roles.id").
        Joins("JOIN permissions ON permissions.id = role_permissions.permission_id").
        Where("users.id = ? AND permissions.name = ?", userID, permission).
        Count(&count).Error

    return count > 0, err
}