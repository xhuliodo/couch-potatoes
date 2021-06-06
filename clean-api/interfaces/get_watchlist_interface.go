package interfaces

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/xhuliodo/couch-potatoes/clean-api/application"
	common_http "github.com/xhuliodo/couch-potatoes/clean-api/common/http"
)

type getWatchlistResource struct {
	getWatchlistService application.GetWatchlistService
}

type watchlistView struct {
	Movie     movieView `json:"movie"`
	TimeAdded int64     `json:"timeAdded"`
}

func (gwr getWatchlistResource) GetWatchlist(w http.ResponseWriter, r *http.Request) {
	userId := getUserId(r)
	limit := getLimit(r)
	skip := getSkip(r)

	watchlist, err := gwr.getWatchlistService.GetWatchlist(userId, limit, skip)
	if err != nil {
		_ = render.Render(w, r, common_http.DetermineErr(err))
		return
	}
	view := []watchlistView{}
	for _, w := range watchlist {
		view = append(view, watchlistView{
			Movie: movieView{
				Id:           w.Id,
				Title:        w.Title,
				ReleaseYear:  w.ReleaseYear,
				MoreInfoLink: w.MoreInfoLink,
			},
			TimeAdded: w.TimeAdded,
		})
	}

	render.Render(w, r, common_http.SendPayload(view))
}
