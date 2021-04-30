package interfaces

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/xhuliodo/couch-potatoes/clean-api/application"
	common_http "github.com/xhuliodo/couch-potatoes/clean-api/common/http"
	"github.com/xhuliodo/couch-potatoes/clean-api/domain"
)

// type movieView struct {
// 	Id           string      `json:"id"`
// 	Title        string      `json:"title"`
// 	ReleaseYear  int         `json:"releaseYear"`
// 	Poster       string      `json:"poster"`
// 	MoreInfoLink string      `json:"moreInfoLink"`
// 	Genres       []genreView `json:"genres"`
// }

type genreView struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type moviesResource struct {
	movieRepo domain.MovieRepo
}

func (mr moviesResource) GetAllGenres(w http.ResponseWriter, r *http.Request) {
	genres, err := mr.movieRepo.GetAllGenres()

	if err != nil {
		_ = render.Render(w, r, common_http.ErrInternal(err))
		return
	}
	view := []genreView{}
	for _, genre := range genres {
		view = append(view, genreView{
			Id:   genre.Id,
			Name: genre.Name,
		})
	}
	render.Respond(w, r, view)
}

type userResource struct {
	userRepo application.UserRepo
}

func (ur userResource) SaveGenrePreferences(w http.ResponseWriter, r *http.Request) {
	err := errors.New("lalala")
	_ = render.Render(w, r, common_http.ErrInternal(err))
}

func AddRoutes(router *chi.Mux, movieRepo domain.MovieRepo, userRepo application.UserRepo) {
	movieResource := moviesResource{movieRepo}
	router.Get("/genres", movieResource.GetAllGenres)

	userResource := userResource{userRepo}
	router.Post("/user/genre-preferences", userResource.SaveGenrePreferences)
}
