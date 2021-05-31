package interfaces

import (
	"fmt"
	"net/http"
	"strconv"

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

	watchlist, err := gwr.getWatchlistService.GetWatchlist(userId, limit, skip)
	if err != nil {
		_ = render.Render(w, r, common_http.ErrInternal(err))
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
	render.Respond(w, r, view)
}
