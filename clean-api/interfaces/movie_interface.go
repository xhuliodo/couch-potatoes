package interfaces

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

type movieView struct {
	Id           string      `json:"id"`
	Title        string      `json:"title"`
	ReleaseYear  int         `json:"releaseYear"`
	Poster       string      `json:"poster"`
	MoreInfoLink string      `json:"moreInfoLink"`
	Genres       []genreView `json:"genres"`
}

type genreView struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func GetAllGenres (w http.ResponseWriter, r *http.Request){
	genres, err:=
}

func AddRoutes(router *chi.Mux)
