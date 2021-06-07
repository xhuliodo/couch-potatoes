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
	// swagger:route POST /watchlist/{movieId} watchlist addToWatchlist
	// add a new movie to a user's watchlist
	//
	// responses:
	// 		201: resourceCreatedView
	// 	404: errorResponse
	// 	503: errorResponse

	userId := getUserId(r)
	movieId := chi.URLParam(r, "movieId")

	if errStack := atwr.addToWatchlistService.AddToWatchlist(userId, movieId); errStack != nil {
		_ = render.Render(w, r, common_http.DetermineErr(errStack))
		return
	}

	successMsg := fmt.Sprintf("you just added to watchlist the movie with id: %s", movieId)
	render.Render(w, r, common_http.ResourceCreated(successMsg))
}
