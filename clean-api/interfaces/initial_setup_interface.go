package interfaces

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/render"
	"github.com/pkg/errors"
	"github.com/xhuliodo/couch-potatoes/clean-api/application"
	common_http "github.com/xhuliodo/couch-potatoes/clean-api/common/http"
)

type setupResource struct {
	setupService application.SetupService
}

type genreView struct {
	Id   string `json:"genreId"`
	Name string `json:"name"`
}

func (sr setupResource) GetAllGenres(w http.ResponseWriter, r *http.Request) {
	genres, err := sr.setupService.GetAllGenres()
	if err != nil {
		_ = render.Render(w, r, common_http.DetermineErr(err))
		return
	}

	view := []genreView{}
	for _, genre := range genres {
		view = append(view, genreView{
			Id:   genre.Id,
			Name: genre.Name,
		})
	}

	render.Render(w, r, common_http.SendPayload(view))
}

type inputSaveGenrePref struct {
	InputGenresUuid []string `json:"genres"`
}

func (sr setupResource) SaveGenrePreferences(w http.ResponseWriter, r *http.Request) {
	var inputSaveGenrePref inputSaveGenrePref
	if err := json.NewDecoder(r.Body).Decode(&inputSaveGenrePref); err != nil || len(inputSaveGenrePref.InputGenresUuid) == 0 {
		cause := errors.New("bad_request")
		errStack := errors.Wrap(cause, "please send genre preferences in the required json format")
		_ = render.Render(w, r, common_http.DetermineErr(errStack))
		return
	}

	userId := getUserId(r)

	if err := sr.setupService.SaveGenrePreferences(userId, inputSaveGenrePref.InputGenresUuid); err != nil {
		_ = render.Render(w, r, common_http.DetermineErr(err))
		return
	}

	render.Render(w, r, common_http.ResourceCreated("genre preferences have been saved"))
}

type SetupStepView struct {
	Step     uint   `json:"step"`
	Finished bool   `json:"finished"`
	Message  string `json:"message"`
}

func (sr setupResource) GetUserSetupStep(w http.ResponseWriter, r *http.Request) {
	userId := getUserId(r)

	setupStep, err := sr.setupService.GetSetupStep(userId)
	if err != nil {
		_ = render.Render(w, r, common_http.DetermineErr(err))
		return
	}

	view := SetupStepView{
		Step:     setupStep.Step,
		Finished: setupStep.Finished,
		Message:  setupStep.Message,
	}

	render.Render(w, r, common_http.SendPayload(view))
}
