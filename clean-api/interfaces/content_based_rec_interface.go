package interfaces

import (
	"fmt"
	"net/http"
	"strconv"

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
	userIdInterface := r.Context().Value("userId")
	userId := fmt.Sprintf("%v", userIdInterface)

	var limit uint
	limitUrlQueryParam := r.URL.Query().Get("limit")
	limitU64, err := strconv.ParseUint(limitUrlQueryParam, 10, 32)
	if err != nil {
		limit = 5
	}
	limit = uint(limitU64)

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
