package application

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/xhuliodo/couch-potatoes/clean-api/domain"
)

type SetupService struct {
	movieRepo domain.MovieRepo
	userRepo  UserRepo
}

func NewSetupService(movieRepo domain.MovieRepo, userRepo UserRepo) SetupService {
	return SetupService{movieRepo, userRepo}
}

func (ss SetupService) GetAllGenres() ([]domain.Genre, error) {
	return ss.movieRepo.GetAllGenres()
}

func (ss SetupService) SaveGenrePreferences(userId UserId, genres []uuid.UUID) error {
	user, err := ss.userRepo.ById(userId)
	if err != nil {
		return errors.New("a user with this identifier does not exist")
	}

	realG, err := ss.movieRepo.GetAllGenres()
	if err != nil {
		return errors.New("no genres found")
	}

	genresToAdd := []domain.Genre{}
	for _, genre := range genres {
		g, found := Find(realG, genre)
		if !found {
			_, errMessage := fmt.Printf("the genre with this ID: %s, does not exist", g)
			return errMessage
		}
		genresToAdd = append(genresToAdd, g)
	}

	if err := user.MovieWatcher.GiveGenrePreferences(genresToAdd); err != nil {
		return err
	}

	return nil
}

func Find(slice []domain.Genre, val uuid.UUID) (domain.Genre, bool) {
	for _, item := range slice {
		if item.Id == val {
			return item, true
		}
	}
	return domain.Genre{}, false
}
