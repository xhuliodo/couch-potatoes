package http

import (
	"net/http"

	"github.com/go-chi/render"
)

// @description takes the form of the response depending on the resource
type Payload interface{} //@name data

// A SuccessResponse is a response that is used when a request has a payload
type SuccessResponse struct {
	HTTPStatusCode int     `json:"statusCode"`        // http response status code
	Payload        Payload `json:"data"`    // any kind of success payload
} //@name SuccessResponse

func (e *SuccessResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

// any kind of success payload
func SendPayload(payload Payload) render.Renderer {
	return &SuccessResponse{
		HTTPStatusCode: http.StatusOK,
		Payload:        payload,
	}
}
