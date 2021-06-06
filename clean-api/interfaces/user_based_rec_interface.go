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
}

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
