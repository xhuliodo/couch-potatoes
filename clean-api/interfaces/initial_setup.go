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
} //@name GenreResponse

// @router /genres [get]
// @param authorization header string true "Bearer token"
// @summary get all genres
// @tags movies
// @produce  json
// @success 200 {object} common_http.SuccessResponse{data=genreView} "api response"
// @failure 401 {object} common_http.ErrorResponse "when a request without a valid Bearer token is provided"
// @failure 503 {object} common_http.ErrorResponse "when the api cannot connect to the database"
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
} //@name GenrePreferencesInput

// @router /users/genres [post]
// @param authorization header string true "Bearer token"
// @param genres body inputSaveGenrePref true "an array of genre ids (minimum number of genreId's required is 3)" minItems(3)
// @summary give genre preference for a user
// @tags users
// @accept json
// @produce json
// @success 201 {object} common_http.InfoResponse "api response"
// @failure 400 {object} common_http.ErrorResponse "when either a genre provided does not exist or the minumum number of genre preferences have not been given"
// @failure 401 {object} common_http.ErrorResponse "when a request without a valid Bearer token is provided"
// @failure 404 {object} common_http.ErrorResponse "when the user making the request has not been registered in the database yet"
// @failure 503 {object} common_http.ErrorResponse "when the api cannot connect to the database"
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

	render.Render(w, r, common_http.Info("genre preferences have been saved"))
}

type setupStepView struct {
	Step         uint   `json:"step"`
	Finished     bool   `json:"finished"`
	Message      string `json:"message"`
	RatingsGiven uint   `json:"ratingsGiven"`
} //@name SetupStepResponse

// @router /users/setup [get]
// @param authorization header string true "Bearer token"
// @summary get user's current step in the setup process
// @tags users
// @produce  json
// @success 200 {object} common_http.SuccessResponse{data=setupStepView} "api response"
// @failure 401 {object} common_http.ErrorResponse "when a request without a valid Bearer token is provided"
// @failure 404 {object} common_http.ErrorResponse "when the user making the request has not been registered in the database yet"
// @failure 503 {object} common_http.ErrorResponse "when the api cannot connect to the database"
func (sr setupResource) GetUserSetupStep(w http.ResponseWriter, r *http.Request) {
	userId := getUserId(r)

	setupStep, err := sr.setupService.GetSetupStep(userId)
	if err != nil {
		_ = render.Render(w, r, common_http.DetermineErr(err))
		return
	}

	view := setupStepView{
		Step:         setupStep.Step,
		Finished:     setupStep.Finished,
		Message:      setupStep.Message,
		RatingsGiven: setupStep.RatingsGiven,
	}

	render.Render(w, r, common_http.SendPayload(view))
}
