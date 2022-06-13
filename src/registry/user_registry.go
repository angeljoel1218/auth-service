package registry

import (
	uui "auth-service/src/app/interactor"
	uur "auth-service/src/app/repository"
	uus "auth-service/src/app/service"
	uir "auth-service/src/infrastructure/database/repository"
)

func (r *registry) NewUserSerivice() uus.UserService {
	return uui.NewUserInteractor(r.NewUserRepository())
}

func (r *registry) NewUserRepository() uur.UserRepository {
	return uir.NewUserRepository(r.db)
}
