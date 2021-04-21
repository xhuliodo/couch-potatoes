package interfaces

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/google/uuid"
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
	repo domain.MovieRepo
}

func (mr moviesResource) GetAllGenres(w http.ResponseWriter, r *http.Request) {
	genres, err := mr.repo.GetAllGenres()

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

func AddRoutes(router *chi.Mux, repo domain.MovieRepo) {
	resource := moviesResource{repo}
	router.Get("/genres", resource.GetAllGenres)
}
