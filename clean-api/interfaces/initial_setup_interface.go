package interfaces

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/xhuliodo/couch-potatoes/clean-api/application"
	common_http "github.com/xhuliodo/couch-potatoes/clean-api/common/http"
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
	InputGenresUuid []string `json:"genres"`
}

func (sr setupResource) SaveGenrePreferences(w http.ResponseWriter, r *http.Request) {
	var inputSaveGenrePref inputSaveGenrePref
	if err := json.NewDecoder(r.Body).Decode(&inputSaveGenrePref); err != nil {
		_ = render.Render(w, r, common_http.ErrInternal(err))
		return
	}

	userIdInterface := r.Context().Value("userId")
	userId := fmt.Sprintf("%v", userIdInterface)

	if err := sr.setupService.SaveGenrePreferences(userId, inputSaveGenrePref.InputGenresUuid); err != nil {
		_ = render.Render(w, r, common_http.ErrInternal(err))
		return
	}

	render.Render(w, r, common_http.ResourceCreated("genre preferences of the user have been saved"))
}

type SetupStepView struct {
	Step     uint   `json:"step"`
	Finished bool   `json:"finished"`
	Message  string `json:"message"`
}

func (sr setupResource) GetUserSetupStep(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")
	setupStep, err := sr.setupService.GetSetupStep(userId)

	if err != nil {
		_ = render.Render(w, r, common_http.ErrInternal(err))
		return
	}
	view := SetupStepView{
		Step:     setupStep.Step,
		Finished: setupStep.Finished,
		Message:  setupStep.Message,
	}

	render.Respond(w, r, view)
}
