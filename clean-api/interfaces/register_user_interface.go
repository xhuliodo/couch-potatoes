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

func (rs registerResource) RegisterUser(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")
	isAdmin := getIsAdmin(r)

	if err := rs.registerService.RegisterUser(userId, isAdmin); err != nil {
		_ = render.Render(w, r, common_http.DetermineErr(err))
		return
	}

	render.Render(w, r, common_http.ResourceCreated("user has been successfully created"))
}
