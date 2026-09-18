package repository

import "github.com/vidhyashekar/cloudcart/services/auth-service/internal/model"

type UserRepository interface {
	// Define methods for user repository
	Create(user *model.User) error
	FindByEmail(email string) (*model.User, error)
	FindRoleByName(name string) (*model.Role, error)
	CreateRefreshToken(token *model.RefreshToken) error
	FindRefreshToken(token string) (*model.RefreshToken, error)
}
