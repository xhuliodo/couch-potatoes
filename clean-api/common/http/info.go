package http

import (
	"net/http"

	"github.com/go-chi/render"
)

type InfoResponse struct {
	HTTPStatusCode int    `json:"statusCode"`        // http response status code
	Message        string `json:"message,omitempty"` // a comment when necessary
} //@name InfoResponse

func (e *InfoResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

// when a new rescource (movie,genre,user) was created
func Info(message string) render.Renderer {
	return &InfoResponse{
		HTTPStatusCode: http.StatusCreated,
		Message:        message,
	}
}

// when there are no more recommendations at
// the moment
func NoContent() render.Renderer {
	return &InfoResponse{
		HTTPStatusCode: http.StatusNoContent,
	}
}
