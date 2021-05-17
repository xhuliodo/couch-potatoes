package application

import (
	"errors"

	"github.com/xhuliodo/couch-potatoes/clean-api/domain"
)

type RegisterService struct {
	repo domain.Repository
}

func NewRegisterService(repo domain.Repository) RegisterService {
	return RegisterService{repo}
}

func (rs RegisterService) RegisterUser(userId string, isAdmin bool) error {
	if !isAdmin {
		return errors.New("you are not allowed to register users")
	}

	newUser := domain.User{Id: userId}

	if err := rs.repo.RegisterNewUser(newUser); err != nil {
		return errors.New("the new user could not be registered")
	}

	return nil
}
