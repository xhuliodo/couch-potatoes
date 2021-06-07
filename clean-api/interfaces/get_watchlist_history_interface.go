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
// @param authorization header string true "Bearer token"
// @summary get all movies in a user's watchlist history
// @tags watchlists
// @produce  json
// @success 200 {object} common_http.SuccessResponse{data=watchlistHistoryView} "api response"
// @failure 404 {object} common_http.ErrorResponse "when either a movie of a user does not exist"
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
