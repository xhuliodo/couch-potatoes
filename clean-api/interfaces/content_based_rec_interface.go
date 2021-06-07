package interfaces

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/xhuliodo/couch-potatoes/clean-api/application"
	common_http "github.com/xhuliodo/couch-potatoes/clean-api/common/http"
)

type contentBasedRecResource struct {
	contentBasedRecService application.ContentBasedRecommendationService
}

type contentBasedRecView struct {
	Movie movieView `json:"movie"`
	Score float64   `json:"score"`
}//@name ContentBasedRecommendationResponse

// @router /recommendations/user-based [get]
// @param limit query int true "limit" default(5)
// @param authorization header string true "Bearer token"
// @summary get content based recommendation from previously liked movies
// @tags recommendations
// @produce json
// @success 200 {object} common_http.SuccessResponse{data=contentBasedRecView} "api response"
// @failure 400 {object} common_http.ErrorResponse "when there are no more recommendations to give"
// @failure 503 {object} common_http.ErrorResponse "when the api cannot connect to the database"
func (cbrr contentBasedRecResource) GetContentBasedRecommendation(w http.ResponseWriter, r *http.Request) {
	userId := getUserId(r)
	limit := getLimit(r)

	contentBasedRec, errStack := cbrr.contentBasedRecService.GetContentBasedRecommendation(userId, limit)
	if errStack != nil {
		_ = render.Render(w, r, common_http.DetermineErr(errStack))
		return
	}

	view := []contentBasedRecView{}
	for _, rec := range contentBasedRec {
		view = append(view, contentBasedRecView{
			Movie: movieView{
				Id:           rec.Id,
				Title:        rec.Title,
				ReleaseYear:  rec.ReleaseYear,
				MoreInfoLink: rec.MoreInfoLink,
			},
			Score: rec.Score,
		})
	}

	render.Respond(w, r, common_http.SendPayload(view))
}
