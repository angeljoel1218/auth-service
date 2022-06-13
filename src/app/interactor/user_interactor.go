package interactor

import (
	"auth-service/src/app/repository"
	us "auth-service/src/app/service"
	apierror "auth-service/src/domain/apierrors"
	fixture "auth-service/src/domain/fixtures"
	"auth-service/src/domain/models"
)

type userInteractor struct {
	UserRepository repository.UserRepository
}

func NewUserInteractor(r repository.UserRepository) us.UserService {
	return &userInteractor{
		UserRepository: r,
	}
}

func (ui *userInteractor) SignIn(input *models.User) (*models.User, error) {

	if err := input.ValidateLogin(); err != nil {
		return nil, err
	}

	user, err := ui.UserRepository.FindByEmail(input.Email)

	if err != nil {
		return nil, err
	}

	if !fixture.ComparePasswords(user.Password, input.Password) {
		return nil, models.ErrorUser(apierror.Unauthorized)
	}

	token, err := fixture.CreateToken(*user.ID)
	if err != nil {
		return nil, err
	}

	user.Token = token

	return user, nil
}

func (ui *userInteractor) SignUp(input *models.User) (*models.User, error) {

	if err := input.ValidateCreation(); err != nil {
		return nil, err
	}

	hashed, err := fixture.HashAndSalt(input.Password)

	if err != nil {
		return nil, err
	}

	user, err := ui.UserRepository.FindByEmail(input.Email)

	if err != nil && apierror.HttpStatus(err) != 404 {
		return nil, err
	}

	if user != nil && user.Email == input.Email {
		return nil, models.ErrorUser(apierror.Conflict)
	}

	input.Password = hashed

	us, err := ui.UserRepository.Create(input)

	if err != nil {
		return nil, err
	}

	token, err := fixture.CreateToken(*us.ID)
	if err != nil {
		return nil, err
	}

	us.Token = token

	return user, nil
}
