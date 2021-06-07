package interfaces

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/pkg/errors"
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
	MovieId string `json:"movieId"`
	Rating  int    `json:"rating"`
} //@name RateMovieInput

// @router /users/ratings [post]
// @param authorization header string true "Bearer token"
// @param rating body inputRateMovie true "an object of movieId and user rating for that movie"
// @summary user rates a movie
// @tags users
// @produce  json
// @success 201 {object} ResourceCreatedView "api response"
// @failure 400 {object} common_http.ErrorResponse "when the input provided was not in the specified json format"
// @failure 404 {object} common_http.ErrorResponse "when the movie provided does not exist"
// @failure 503 {object} common_http.ErrorResponse "when the api cannot connect to the database"
func (rs ratingResource) RateMovie(w http.ResponseWriter, r *http.Request) {
	var inputRateMovie inputRateMovie
	if err := json.NewDecoder(r.Body).Decode(&inputRateMovie); err != nil || isInputEmpty(inputRateMovie) {
		cause := errors.New("bad_request")
		errStack := errors.Wrap(cause, "please send movieId and rating in the required json format")
		_ = render.Render(w, r, common_http.DetermineErr(errStack))
		return
	}

	userId := getUserId(r)

	if err := rs.ratingService.RateMovie(userId, inputRateMovie.MovieId, inputRateMovie.Rating); err != nil {
		_ = render.Render(w, r, common_http.DetermineErr(err))
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

func isInputEmpty(input inputRateMovie) bool {
	if input.MovieId == "" || input.Rating == 0 {
		return true
	}
	return false
}
