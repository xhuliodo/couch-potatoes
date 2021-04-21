package use_cases

import "github.com/xhuliodo/couch-potatoes/clean-api/domain"

type User struct {
	isAdmin      bool
	MovieWatcher domain.MovieWatcher
}

type SetupService struct {
	movieRepo domain.MovieRepo
}

func NewSetupService(movieRepo domain.MovieRepo) SetupService {
	return SetupService{movieRepo}
}

func (ss SetupService) GetAllGenres() ([]domain.Genre, error) {
	return ss.movieRepo.GetAllGenres()
}
