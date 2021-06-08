package application

import (
	"github.com/pkg/errors"
	"github.com/xhuliodo/couch-potatoes/clean-api/domain"
)

type SetupService struct {
	repo domain.Repository
}

func NewSetupService(repo domain.Repository) SetupService {
	return SetupService{repo}
}

func (ss SetupService) GetAllGenres() ([]domain.Genre, error) {
	genres, err := ss.repo.GetAllGenres()
	if err != nil {
		errStack := errors.Wrap(err, "could not retrieve all genres")
		return genres, errStack
	}

	return genres, nil
}

func (ss SetupService) SaveGenrePreferences(userId string, genres []string) error {
	user, err := ss.repo.GetUserById(userId)
	if err != nil {
		errStack := errors.Wrap(err, "a user with this identifier does not exist")
		return errStack
	}

	currentGenres, err := ss.repo.GetAllGenres()
	if err != nil {
		errStack := errors.Wrap(err, "could not retrieve all genres")
		return errStack
	}

	genresToAdd := []domain.Genre{}
	for _, genre := range genres {
		g, found := Find(currentGenres, genre)
		if !found {
			cause := errors.New("bad_request")
			return errors.Wrapf(cause, "genre with id %s does not exist", genre)
		}
		genresToAdd = append(genresToAdd, g)
	}

	if err := user.GiveGenrePreferences(genresToAdd); err != nil {
		cause := errors.New("bad_request")
		return errors.Wrap(cause, err.Error())
	}

	if err := ss.repo.SaveGenrePreferences(user.Id, genresToAdd); err != nil {
		errStack := errors.Wrap(err, "could not save genre preferences")
		return errStack
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

type SetupStep struct {
	Step         uint
	Finished     bool
	Message      string
	RatingsGiven uint
}

const (
	leastRequiredGenrePref uint = 3
	leastRequiredRatingNr  uint = 15
)

func (ss SetupService) GetSetupStep(userId string) (
	sp SetupStep, err error,
) {
	if _, err := ss.repo.GetUserById(userId); err != nil {
		errStack := errors.Wrap(err, "a user with this identifier does not exist")
		return sp, errStack
	}

	genrePrefCount, err := ss.repo.GetGenrePreferencesCount(userId)
	if err != nil {
		errStack := errors.Wrap(err, "could not get genre preferences count")
		return sp, errStack
	}

	if genrePrefCount < leastRequiredGenrePref {
		sp = SetupStep{
			Step:     1,
			Finished: false,
			Message:  "user has yet to give genre preferences",
		}
		return sp, nil
	}

	ratedMoviesCount, err := ss.repo.GetUserRatingsCount(userId)
	if err != nil {
		errStack := errors.Wrap(err, "could not get user ratings count")
		return sp, errStack
	}

	if ratedMoviesCount < leastRequiredRatingNr {
		sp = SetupStep{
			Step:         2,
			Finished:     false,
			Message:      "user has not given enough ratings to get to know him yet",
			RatingsGiven: ratedMoviesCount,
		}
		return sp, nil
	}

	sp = SetupStep{
		Step:         3,
		Finished:     true,
		Message:      "user has finished the setup process",
		RatingsGiven: ratedMoviesCount,
	}

	return sp, nil
}
