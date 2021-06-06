package interfaces

import (
	"net/http"

	"github.com/go-chi/render"
	common_http "github.com/xhuliodo/couch-potatoes/clean-api/common/http"
)

type watchlistHistoryView struct {
	Movie     movieView `json:"movie"`
	Rating    float64   `json:"rating"`
	TimeAdded int64     `json:"timeAdded"`
}

func (gwr getWatchlistResource) GetWatchlistHistory(w http.ResponseWriter, r *http.Request) {
	userId := getUserId(r)
	limit := getLimit(r)
	skip := getSkip(r)

	watchlistHistory, err := gwr.getWatchlistService.GetWatchlistHistory(userId, limit, skip)
	if err != nil {
		_ = render.Render(w, r, common_http.DetermineErr(err))
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

	render.Render(w, r, common_http.SendPayload(view))
}
