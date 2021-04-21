package infrastructure

import "github.com/xhuliodo/couch-potatoes/clean-api/domain"

type InMemoryRepository struct {
	movie []domain.Movie
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{[]domain.Movie{}}
}

func (m *InMemoryRepository) GetAllGenres([]domain.Genre, error) {

}
