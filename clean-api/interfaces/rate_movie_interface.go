package interfaces

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/xhuliodo/couch-potatoes/clean-api/application"
	common_http "github.com/xhuliodo/couch-potatoes/clean-api/common/http"
)

type ratingResource struct {
	ratingService application.RatingService
}

func NewRatingResource(ratingService application.RatingService) ratingResource {
	return ratingResource{ratingService}
}

type inputRateMovie struct {
	// UserId  string `json:"userId"`
	MovieId string `json:"movieId"`
	Rating  int    `json:"rating"`
}

func (rs ratingResource) RateMovie(w http.ResponseWriter, r *http.Request) {
	var inputRateMovie inputRateMovie
	if err := json.NewDecoder(r.Body).Decode(&inputRateMovie); err != nil {
		_ = render.Render(w, r, common_http.ErrInternal(err))
		return
	}

	userIdInterface := r.Context().Value("userId")
	userId := fmt.Sprintf("%v", userIdInterface)

	if err := rs.ratingService.RateMovie(userId, inputRateMovie.MovieId, inputRateMovie.Rating); err != nil {
		_ = render.Render(w, r, common_http.ErrInternal(err))
		return
	}

	var ratingMeasure string
	if inputRateMovie.Rating == 0 {
		ratingMeasure = "disliked"
	} else {
		ratingMeasure = "liked"
	}

	successMsg := fmt.Sprintf("you just %s the movie with id: %s", ratingMeasure, inputRateMovie.MovieId)
	render.Render(w, r, common_http.ResourceCreated(successMsg))
}
