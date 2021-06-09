package application

import (
	"github.com/pkg/errors"
	"github.com/xhuliodo/couch-potatoes/clean-api/domain"
)

type HealthcheckService struct {
	repo domain.Repository
}

func NewHealthcheckService(repo domain.Repository) HealthcheckService {
	return HealthcheckService{repo}
}

func (hs HealthcheckService) Heathcheck() error {
	err := hs.repo.Healthcheck()
	if err != nil {
		errStack := errors.Wrap(err, "cannot connect to database")
		return errStack
	}
	
	return nil
}
