package interfaces

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/render"
	"github.com/xhuliodo/couch-potatoes/clean-api/application"
	common_http "github.com/xhuliodo/couch-potatoes/clean-api/common/http"
)

type registerResource struct {
	registerService application.RegisterService
}

type inputRegisterUser struct {
	UserId string `json:"userId"`
}

func (rs registerResource) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var inputRegisterUser inputRegisterUser
	if err := json.NewDecoder(r.Body).Decode(&inputRegisterUser); err != nil {
		_ = render.Render(w, r, common_http.ErrInternal(err))
		return
	}

	isAdminInterface := r.Context().Value("isAdmin")
	isAdmin := isAdminInterface.(bool)

	if err := rs.registerService.RegisterUser(inputRegisterUser.UserId, isAdmin); err != nil {
		_ = render.Render(w, r, common_http.ErrInternal(err))
		return
	}

	render.Render(w, r, common_http.ResourceCreated("user has been successfully created"))
}
