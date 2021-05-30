package interfaces

import (
	"fmt"
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
	userIdInterface := r.Context().Value("userId")
	userId := fmt.Sprintf("%v", userIdInterface)

	movieId := chi.URLParam(r, "movieId")

	if err := rfwr.removeFromWatchlistService.RemoveFromWatchlist(userId, movieId); err != nil {
		_ = render.Render(w, r, common_http.ErrInternal(err))
		return
	}

	successMsg := fmt.Sprintf("you just removed from watchlist the movie with id: %s", movieId)
	render.Render(w, r, common_http.ResourceCreated(successMsg))
}
