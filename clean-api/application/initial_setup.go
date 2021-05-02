package application

import (
	"errors"

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

func (ss SetupService) SaveGenrePreferences(userId string, genres []uuid.UUID) error {
	user, err := ss.userRepo.GetUserById(userId)
	if err != nil {
		return errors.New("a user with this identifier does not exist")
	}

	currentGenres, err := ss.movieRepo.GetAllGenres()
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

	if err := user.MovieWatcher.GiveGenrePreferences(genresToAdd); err != nil {
		return err
	}

	if err := ss.userRepo.SaveGenrePreferences(user.Id, genresToAdd); err != nil {
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
