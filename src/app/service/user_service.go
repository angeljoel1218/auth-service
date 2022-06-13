package service

import "auth-service/src/domain/models"

type UserService interface {
	SignIn(user *models.User) (*models.User, error)
	SignUp(user *models.User) (*models.User, error)
}
