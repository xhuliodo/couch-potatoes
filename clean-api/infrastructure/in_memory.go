package infrastructure

import (
	"github.com/google/uuid"
	"github.com/xhuliodo/couch-potatoes/clean-api/domain"
)

type InMemoryRepository struct {
	movies []domain.Movie
	genres []domain.Genre
}

func NewInMemoryRepository() *InMemoryRepository {
	uuid1 := uuid.MustParse("1ae88fa0-a355-11eb-bcbc-0242ac130002")
	uuid2 := uuid.MustParse("1ae891d0-a355-11eb-bcbc-0242ac130002")
	return &InMemoryRepository{[]domain.Movie{}, []domain.Genre{
		{Id: uuid1, Name: "Adventure"},
		{Id: uuid2, Name: "Romance"},
	}}
}

func (imr *InMemoryRepository) GetAllGenres() ([]domain.Genre, error) {
	return imr.genres, nil
}
