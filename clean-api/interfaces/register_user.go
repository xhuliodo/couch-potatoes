package interfaces

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/xhuliodo/couch-potatoes/clean-api/application"
	common_http "github.com/xhuliodo/couch-potatoes/clean-api/common/http"
)

type registerResource struct {
	registerService application.RegisterService
}

// @router /users/{userId} [post]
// @param userId path string true "the id of the new user"
// @param authorization header string true "Bearer token"
// @summary registering a user to the app, it can be done only through users with the role admin
// @tags users
// @produce json
// @success 201 {object} common_http.InfoResponse "api response"
// @failure 401 {object} common_http.ErrorResponse "when a request without a valid Bearer token is provided"
// @failure 404 {object} common_http.ErrorResponse "when the user making the request has not been registered in the database yet"
// @failure 503 {object} common_http.ErrorResponse "when the api cannot connect to the database"
func (rs registerResource) RegisterUser(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")
	isAdmin := getIsAdmin(r)

	if err := rs.registerService.RegisterUser(userId, isAdmin); err != nil {
		_ = render.Render(w, r, common_http.DetermineErr(err))
		return
	}

	render.Render(w, r, common_http.Info("user has been successfully created"))
}
