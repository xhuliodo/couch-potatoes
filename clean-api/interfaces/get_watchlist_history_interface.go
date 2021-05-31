package interfaces

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/render"
	common_http "github.com/xhuliodo/couch-potatoes/clean-api/common/http"
)

type watchlistHistoryView struct {
	Movie     movieView `json:"movie"`
	Rating    float64   `json:"rating"`
	TimeAdded int64     `json:"timeAdded"`
}

func (gwr getWatchlistResource) GetWatchlistHistory(w http.ResponseWriter, r *http.Request) {
	userIdInterface := r.Context().Value("userId")
	userId := fmt.Sprintf("%v", userIdInterface)

	var limit uint
	limitUrlQueryParam := r.URL.Query().Get("limit")
	limitU64, err := strconv.ParseUint(limitUrlQueryParam, 10, 32)
	if err != nil {
		limit = 5
	}
	limit = uint(limitU64)

	var skip uint
	skipUrlQueryParam := r.URL.Query().Get("skip")
	skipU64, err := strconv.ParseUint(skipUrlQueryParam, 10, 32)
	if err != nil {
		skip = 5
	}
	skip = uint(skipU64)

	watchlistHistory, err := gwr.getWatchlistService.GetWatchlistHistory(userId, limit, skip)
	if err != nil {
		_ = render.Render(w, r, common_http.ErrInternal(err))
		return
	}
	view := []watchlistHistoryView{}
	for _, w := range watchlistHistory {
		view = append(view, watchlistHistoryView{
			Movie: movieView{
				Id:           w.Id,
				Title:        w.Title,
				ReleaseYear:  w.ReleaseYear,
				MoreInfoLink: w.MoreInfoLink,
			},
			Rating:    w.Rating,
			TimeAdded: w.TimeAdded,
		})
	}
	render.Respond(w, r, view)
}
