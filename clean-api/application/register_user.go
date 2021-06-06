package application

import (
	"github.com/pkg/errors"
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
		cause := errors.New("forbidden")
		return errors.Wrap(cause, "you are not allowed to register users")
	}

	if err := rs.repo.RegisterNewUser(userId); err != nil {
		return errors.Wrap(err, "new user could not be registered")
	}

	return nil
}
