package repository

import (
	apierror "auth-service/src/domain/apierrors"
	"auth-service/src/domain/models"
	"auth-service/src/usecase/repository"
	"errors"

	"gorm.io/gorm"
)

type userRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepository{
		DB: db,
	}
}

//function create user
func (r *userRepository) Create(user *models.User) (*models.User, error) {
	if err := r.DB.Create(user).Error; !errors.Is(err, nil) {
		return nil, err
	}

	return user, nil
}

// search user form email
func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	findUser := &models.User{}
	r.DB.Find(findUser, "email = ?", email)

	if findUser.ID == nil {
		return nil, models.ErrorUser(apierror.NotFound)
	}

	return findUser, nil
}
