package application

import (
	"errors"
)

type RegisterService struct {
	userRepo UserRepo
}

func NewRegisterService(userRepo UserRepo) RegisterService {
	return RegisterService{userRepo}
}

func (rs RegisterService) RegisterUser(userId string, isAdmin bool) error {
	if !isAdmin {
		return errors.New("you are not allowed to register users")
	}

	newUser := User{Id: userId}

	if err := rs.userRepo.RegisterNewUser(newUser); err != nil {
		return errors.New("the new user could not be registered")
	}

	return nil
}
