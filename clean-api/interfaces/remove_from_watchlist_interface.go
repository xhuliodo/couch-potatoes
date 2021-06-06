package interfaces

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/xhuliodo/couch-potatoes/clean-api/application"
	common_http "github.com/xhuliodo/couch-potatoes/clean-api/common/http"
)

type removeFromWatchlistResource struct {
	removeFromWatchlistService application.RemoveFromWatchlistService
}

func NewRemoveFromWatchlistResource(removeFromWatchlistService application.RemoveFromWatchlistService) removeFromWatchlistResource {
	return removeFromWatchlistResource{removeFromWatchlistService}
}

func (rfwr removeFromWatchlistResource) RemoveFromWatchlist(w http.ResponseWriter, r *http.Request) {
	userId := getUserId(r)
	movieId := chi.URLParam(r, "movieId")

	if errStack := rfwr.removeFromWatchlistService.RemoveFromWatchlist(userId, movieId); errStack != nil {
		_ = render.Render(w, r, common_http.DetermineErr(errStack))
		return
	}

	render.Render(w, r, common_http.NoContent())
}
