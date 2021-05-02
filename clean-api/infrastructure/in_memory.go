package infrastructure

import (
	"errors"

	"github.com/google/uuid"
	"github.com/xhuliodo/couch-potatoes/clean-api/application"
	"github.com/xhuliodo/couch-potatoes/clean-api/domain"
)

type InMemoryRepository struct {
	movies []domain.Movie
	genres []domain.Genre
	users  []application.User
}

func NewInMemoryRepository() *InMemoryRepository {
	uuid1 := uuid.MustParse("c4f88090-9166-4ebf-920b-ff9a34872b84")
	uuid2 := uuid.MustParse("acffe5b6-d327-43f6-b5ca-0a86f6780629")
	uuid3 := uuid.MustParse("35ee6205-1b06-4eff-8efd-f396ede8a52e")
	uuid4 := uuid.MustParse("3e70ce4e-ae21-463e-bb92-575204f83cd0")
	userUuid := "d3e61774-206e-4c8b-82d4-7f46fb18a3ad"
	return &InMemoryRepository{
		movies: []domain.Movie{},
		genres: []domain.Genre{
			{Id: uuid1, Name: "Adventure"},
			{Id: uuid2, Name: "Romance"},
			{Id: uuid3, Name: "Comedy"},
			{Id: uuid4, Name: "Thriller"},
		},
		users: []application.User{
			{
				Id:      userUuid,
				IsAdmin: false,
				MovieWatcher: domain.MovieWatcher{
					Id:             domain.MovieWatcherID(userUuid),
					Name:           "Chulio",
					FavoriteGenres: []domain.Genre{},
					RatedMovies:    []domain.RatedMovie{},
					Watchlist:      []domain.Movie{},
				},
			},
		},
	}

}

func (imr *InMemoryRepository) GetAllGenres() ([]domain.Genre, error) {
	return imr.genres, nil
}

func (imr *InMemoryRepository) GetUserById(userId string) (application.User, error) {
	u, found := Find(imr.users, userId)
	if !found {
		return application.User{}, errors.New("nope")
	}
	return u, nil
}

func (imr *InMemoryRepository) SaveGenrePreferences(userId string, genres []domain.Genre) error {
	for _, user := range imr.users {
		if user.Id == userId {
			user.MovieWatcher.FavoriteGenres = append(user.MovieWatcher.FavoriteGenres, genres...)
		}
	}
	return nil
}

func Find(slice []application.User, val string) (application.User, bool) {
	for _, item := range slice {
		if item.Id == val {
			return item, true
		}
	}
	return application.User{}, false
}
