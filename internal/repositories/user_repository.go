package repositories

import (
	"stocet-core-service/internal/database"
	"stocet-core-service/internal/models"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) Create(user *models.User) error {
	return database.DB.Create(user).Error
}

func (r *UserRepository) GetByID(id string) (*models.User, error) {
	var user models.User
	err := database.DB.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	err := database.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Update(user *models.User) error {
	return database.DB.Save(user).Error
}

func (r *UserRepository) Delete(id string) error {
	return database.DB.Delete(&models.User{}, "id = ?", id).Error
}

func (r *UserRepository) List(limit, offset int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	if err := database.DB.Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := database.DB.Limit(limit).Offset(offset).Find(&users).Error
	return users, total, err
}

func (r *UserRepository) SetVerifiedByEmail(email string) error {
	return database.DB.Model(&models.User{}).
		Where("email = ?", email).
		Update("is_verified", true).Error
}
