package interfaces

import (
	"net/http"

	"github.com/xhuliodo/couch-potatoes/clean-api/application")


type userResource struct {
	repo application.UserRepo
}

func (ur userResource) SaveGenrePreferences(w http.ResponseWriter, r *http.Request) {
	
}