package application

import "github.com/xhuliodo/couch-potatoes/clean-api/domain"

type PopularMovieService struct {
	repo domain.Repository
}

func NewPopularMovieService(repo domain.Repository) PopularMovieService {
	return PopularMovieService{repo}
}
