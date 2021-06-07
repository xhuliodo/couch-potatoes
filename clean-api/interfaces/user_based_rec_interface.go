package interfaces

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/xhuliodo/couch-potatoes/clean-api/application"
	common_http "github.com/xhuliodo/couch-potatoes/clean-api/common/http"
)

type userBasedRecResource struct {
	userBasedRecService application.UserBasedRecommendationService
}

type userBasedRecView struct {
	Movie movieView `json:"movie"`
	Score float64   `json:"score"`
}//@name UserBasedRecommendationResponse

// @router /recommendations/content-based [get]
// @param limit query int true "limit" default(5)
// @param authorization header string true "Bearer token"
// @summary get users based recommendation from similar users
// @tags recommendations
// @produce json
// @success 200 {object} common_http.SuccessResponse{data=userBasedRecView} "api response"
// @failure 400 {object} common_http.ErrorResponse "when there are no more recommendations to give"
// @failure 503 {object} common_http.ErrorResponse "when the api cannot connect to the database"
func (ubrr userBasedRecResource) GetUserBasedRecommendation(w http.ResponseWriter, r *http.Request) {
	userId := getUserId(r)
	limit := getLimit(r)

	userBasedRec, err := ubrr.userBasedRecService.GetUserBasedRecommendation(userId, limit)
	if err != nil {
		render.Render(w, r, common_http.DetermineErr(err))
		return
	}

	view := []userBasedRecView{}
	for _, rec := range userBasedRec {
		view = append(view, userBasedRecView{
			Movie: movieView{
				Id:           string(rec.Id),
				Title:        rec.Title,
				ReleaseYear:  rec.ReleaseYear,
				MoreInfoLink: rec.MoreInfoLink,
			},
			Score: rec.Score,
		})
	}

	render.Render(w, r, common_http.SendPayload(view))
}
