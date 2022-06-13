package repository

import "auth-service/src/domain/models"

type UserRepository interface {
	FindByEmail(email string) (*models.User, error)
	Create(user *models.User) (*models.User, error)
}
