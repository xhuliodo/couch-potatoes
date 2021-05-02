package http

import (
	"net/http"

	"github.com/go-chi/render"
)

type SuccessResponse struct {
	Message        string `json:"message"` // low-level runtime error
	HTTPStatusCode int    `json:"-"` // http response status code

	AppCode int64 `json:"code,omitempty"` // application-specific error code
}

func (e *SuccessResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ResourceCreated(message string) render.Renderer {
	return &SuccessResponse{
		Message:        message,
		HTTPStatusCode: http.StatusCreated,
	}
}
