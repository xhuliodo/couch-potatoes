package interfaces

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
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
	Id   string `json:"genreId"`
	Name string `json:"name"`
}

type setupResource struct {
	setupService application.SetupService
}

func (sr setupResource) GetAllGenres(w http.ResponseWriter, r *http.Request) {
	genres, err := sr.setupService.GetAllGenres()

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

type inputSaveGenrePref struct {
	UserId          string   `json:"userId"`
	InputGenresUuid []string `json:"genres"`
}

func (sr setupResource) SaveGenrePreferences(w http.ResponseWriter, r *http.Request) {
	var saveGenrePref inputSaveGenrePref
	if err := json.NewDecoder(r.Body).Decode(&saveGenrePref); err != nil {
		_ = render.Render(w, r, common_http.ErrInternal(err))
		return
	}

	if err := sr.setupService.SaveGenrePreferences(saveGenrePref.UserId, saveGenrePref.InputGenresUuid); err != nil {
		_ = render.Render(w, r, common_http.ErrInternal(err))
		return
	}

	render.Render(w, r, common_http.ResourceCreated("genre preferences of the user have been saved"))
}

type userView struct {
	UserId string `json:"userId"`
}

func (sr setupResource) GetUserById(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")
	user, err := sr.setupService.GetUserById(userId)

	if err != nil {
		_ = render.Render(w, r, common_http.ErrInternal(err))
		return
	}
	view := userView{user.Id}

	render.Respond(w, r, view)
}

func AddRoutes(router *chi.Mux, movieRepo domain.MovieRepo, userRepo application.UserRepo) {
	setupService := application.NewSetupService(movieRepo, userRepo)
	setupResource := setupResource{setupService}

	router.Get("/genres", setupResource.GetAllGenres)
	router.Get("/users/{userId}", setupResource.GetUserById)
	router.Post("/users/genres", setupResource.SaveGenrePreferences)
}
