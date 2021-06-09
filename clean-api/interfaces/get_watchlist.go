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
} //@name WatchlistResponse

// @router /watchlist [get]
// @param skip query int false "skip" default(0)
// @param limit query int false "limit" default(5)
// @param authorization header string true "Bearer token"
// @summary get all movies in a user's watchlist
// @tags watchlists
// @produce  json
// @success 200 {object} common_http.SuccessResponse{data=watchlistView} "api response"
// @failure 400 {object} common_http.ErrorResponse "when the skip query param gets too big"
// @failure 401 {object} common_http.ErrorResponse "when a request without a valid Bearer token is provided"
// @failure 404 {object} common_http.ErrorResponse "when there are no movies in the user's watchlist"
// @failure 503 {object} common_http.ErrorResponse "when the api cannot connect to the database"
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
