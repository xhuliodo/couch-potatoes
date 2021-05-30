package interfaces

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/xhuliodo/couch-potatoes/clean-api/application"
	common_http "github.com/xhuliodo/couch-potatoes/clean-api/common/http"
)

type addToWatchlistResource struct {
	addToWatchlistService application.AddToWatchlistService
}

func NewAddToWatchlistResource(addToWatchlistService application.AddToWatchlistService) addToWatchlistResource {
	return addToWatchlistResource{addToWatchlistService}
}

func (atwr addToWatchlistResource) AddToWatchlist(w http.ResponseWriter, r *http.Request) {
	userIdInterface := r.Context().Value("userId")
	userId := fmt.Sprintf("%v", userIdInterface)

	movieId := chi.URLParam(r, "movieId")

	if err := atwr.addToWatchlistService.AddToWatchlist(userId, movieId); err != nil {
		_ = render.Render(w, r, common_http.ErrInternal(err))
		return
	}

	successMsg := fmt.Sprintf("you just added to watchlist the movie with id: %s", movieId)
	render.Render(w, r, common_http.ResourceCreated(successMsg))
}
