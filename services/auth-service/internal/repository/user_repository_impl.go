package repository

import (
	"github.com/vidhyashekar/cloudcart/services/auth-service/internal/model"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User

	err := r.db.
		Preload("Role").
		Where("email = ?", email).
		First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) FindRoleByName(name string) (*model.Role, error) {
	var role model.Role

	err := r.db.
		Where("name = ?", name).
		First(&role).Error

	if err != nil {
		return nil, err
	}

	return &role, nil
}

func (r *userRepository) CreateRefreshToken(
	token *model.RefreshToken,
) error {
	return r.db.Create(token).Error
}

func (r *userRepository) FindRefreshToken(
	token string,
) (*model.RefreshToken, error) {

	var refreshToken model.RefreshToken

	err := r.db.
		Preload("User").
		Preload("User.Role").
		Where("token = ?", token).
		First(&refreshToken).Error

	if err != nil {
		return nil, err
	}

	return &refreshToken, nil
}
