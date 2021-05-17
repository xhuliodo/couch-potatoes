package application

import (
	"errors"

	"github.com/xhuliodo/couch-potatoes/clean-api/domain"
)

type SetupService struct {
	repo domain.Repository
}

func NewSetupService(repo domain.Repository) SetupService {
	return SetupService{repo}
}

func (ss SetupService) GetAllGenres() ([]domain.Genre, error) {
	return ss.repo.GetAllGenres()
}

// func (ss SetupService) GetUserById(userId string) (User, error) {
// 	return ss.userRepo.GetUserById(userId)
// }

func (ss SetupService) SaveGenrePreferences(userId string, genres []string) error {
	user, err := ss.repo.GetUserById(userId)
	if err != nil {
		return errors.New("a user with this identifier does not exist")
	}

	currentGenres, err := ss.repo.GetAllGenres()
	if err != nil {
		return errors.New("no genres found")
	}

	genresToAdd := []domain.Genre{}
	for _, genre := range genres {
		g, found := Find(currentGenres, genre)
		if !found {
			return errors.New("one of the genres id's supplied is incorrect")
		}
		genresToAdd = append(genresToAdd, g)
	}

	if err := user.GiveGenrePreferences(genresToAdd); err != nil {
		return err
	}

	if err := ss.repo.SaveGenrePreferences(user.Id, genresToAdd); err != nil {
		return err
	}

	return nil
}

func Find(slice []domain.Genre, val string) (domain.Genre, bool) {
	for _, item := range slice {
		if item.Id == val {
			return item, true
		}
	}
	return domain.Genre{}, false
}
