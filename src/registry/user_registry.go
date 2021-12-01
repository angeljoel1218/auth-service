package registry

import (
	uip "auth-service/src/interface/presenter"
	uir "auth-service/src/interface/repository"
	uui "auth-service/src/usecase/interactor"
	uup "auth-service/src/usecase/presenter"
	uur "auth-service/src/usecase/repository"
)

func (r *registry) NewUserInteractor() uui.UserInteractor {
	return uui.NewUserInteractor(r.NewUserRepository(), r.NewUserPresenter())
}

func (r *registry) NewUserRepository() uur.UserRepository {
	return uir.NewUserRepository(r.db)
}

func (r *registry) NewUserPresenter() uup.UserPresenter {
	return uip.NewUserPresenter()
}
