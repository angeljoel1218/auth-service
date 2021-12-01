package presenter

import (
	"auth-service/src/domain/models"
	"auth-service/src/usecase/presenter"
)

type userPresenter struct{}

func NewUserPresenter() presenter.UserPresenter {
	return &userPresenter{}
}

func (up *userPresenter) UserResponse(user *models.User) *models.User {
	user.Password = ""
	return user
}
