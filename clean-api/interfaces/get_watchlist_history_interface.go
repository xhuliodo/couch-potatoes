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
}//@name WatchlistHistoryResponse

// @router /watchlist/history [get]
// @param skip query int false "skip" default(0)
// @param limit query int false "limit" default(5)
// @param authorization header string true "Bearer token"
// @summary get all movies in a user's watchlist history
// @tags watchlists
// @produce  json
// @success 200 {object} common_http.SuccessResponse{data=watchlistHistoryView} "api response"
// @failure 400 {object} common_http.ErrorResponse "when the skip query param gets too big"
// @failure 401 {object} common_http.ErrorResponse "when a request without a valid Bearer token is provided"
// @failure 404 {object} common_http.ErrorResponse "when there are no movies in the user's watchlist"
// @failure 503 {object} common_http.ErrorResponse "when the api cannot connect to the database"
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
