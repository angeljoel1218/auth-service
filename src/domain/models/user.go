package models

import (
	apierror "auth-service/src/domain/apierrors"
	"net/mail"
	"time"
)

// User is data user
type User struct {
	ID        *int64     `form:"id" json:"id"`
	Name      string     `form:"name" json:"name"`
	Lastname  string     `form:"lastname" json:"lastname"`
	Email     string     `gorm:"unique" form:"email" json:"email"`
	Password  string     `form:"password" json:"password"`
	Phone     string     `form:"phone" json:"phone,omitempty"`
	CreatedAt time.Time  `json:"created_at,omitempty"`
	UpdatedAt time.Time  `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	Token     string     `gorm:"-" form:"token" json:"token,omitempty"`
}

func (user *User) ValidateCreation() *apierror.ErrorApi {
	if len(user.Name) == 0 {
		return ErrorUser(apierror.BadRequest, "The name field is required")
	}
	if len(user.Lastname) == 0 {
		return ErrorUser(apierror.BadRequest, "The name field is required")
	}
	if len(user.Password) < 6 {
		return ErrorUser(apierror.BadRequest, "The password field must be greater than or equal to 6 characters")
	}
	if _, err := mail.ParseAddress(user.Email); err != nil {
		return ErrorUser(apierror.BadRequest, "Invalid e-mail")
	}
	return nil
}

func (user *User) ValidateLogin() *apierror.ErrorApi {
	if len(user.Password) < 6 {
		return ErrorUser(apierror.BadRequest, "The password field must be greater than or equal to 6 characters")
	}
	if _, err := mail.ParseAddress(user.Email); err != nil {
		return ErrorUser(apierror.BadRequest, "Invalid e-mail")
	}
	return nil
}

func ErrorUser(args ...interface{}) *apierror.ErrorApi {
	switch args[0] {
	case apierror.Conflict:
		return apierror.NewErrorApi(apierror.Conflict, "Username already exists")
	case apierror.NotFound:
		return apierror.NewErrorApi(apierror.NotFound, "Username does not exist")
	case apierror.Unauthorized:
		return apierror.NewErrorApi(apierror.Unauthorized, "Invalid password")
	case apierror.BadRequest:
		return apierror.NewErrorApi(apierror.BadRequest, args[1].(string))
	case apierror.AppError:
		return apierror.NewErrorLogic(apierror.AppError, args[1].(int), args[2].(string))
	default:
		return apierror.NewErrorApi(apierror.InternalServer, "An error occurred, try again")
	}
}
