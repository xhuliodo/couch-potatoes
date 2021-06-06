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
}

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
