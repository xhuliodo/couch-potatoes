package interfaces

import (
	"fmt"
	"net/http"
	"strconv"

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
	userIdInterface := r.Context().Value("userId")
	userId := fmt.Sprintf("%v", userIdInterface)

	var limit uint
	limitUrlQueryParam := r.URL.Query().Get("limit")
	limitU64, err := strconv.ParseUint(limitUrlQueryParam, 10, 32)
	if err != nil {
		limit = 5
	}
	limit = uint(limitU64)

	userBasedRec, err := ubrr.userBasedRecService.GetUserBasedRecommendation(userId, limit)
	if err != nil {
		_ = render.Render(w, r, common_http.ErrInternal(err))
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

	render.Respond(w, r, view)
}
