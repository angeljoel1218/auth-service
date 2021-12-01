package presenter

import "auth-service/src/domain/models"

type UserPresenter interface {
	UserResponse(user *models.User) *models.User
}
