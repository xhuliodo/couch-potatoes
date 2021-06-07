package http

import (
	"net/http"

	"github.com/go-chi/render"
)

// A ErrorResponse is a response that is used when a request cannot be fulfilled, it ranges from the 4XX to 5XX
type ErrorResponse struct {
	ErrStack       error  `json:"-"`               // low-level runtime error
	HTTPStatusCode int    `json:"statusCode"`      // http response status code
	ErrorText      string `json:"error,omitempty"` // application-level error message
} //@name ErrorResponse

func (e *ErrorResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

// ERROR CODES 4XX

// if a user gives wrong arguments to the api
// ERROR_CAUSE: bad_request
func ErrBadRequest(errText string, stackTrace error) render.Renderer {
	return &ErrorResponse{
		ErrStack:       stackTrace,
		HTTPStatusCode: http.StatusBadRequest,
		ErrorText:      errText,
	}
}

// when a user is trying to access api without
// presenting a token
// ERROR_CAUSE: not_authenticated
func ErrUnauthorized(errText string, stackTrace error) render.Renderer {
	return &ErrorResponse{
		ErrStack:       stackTrace,
		HTTPStatusCode: http.StatusUnauthorized,
		ErrorText:      errText,
	}
}

// when a user without the admin role is trying
// to register a new user
// OR
// when a user is trying to get recommendations
// without finishing up the setup process
// ERROR_CAUSE: not_authorized
func ErrForbidden(errText string, stackTrace error) render.Renderer {
	return &ErrorResponse{
		ErrStack:       stackTrace,
		HTTPStatusCode: http.StatusForbidden,
		ErrorText:      errText,
	}
}

// A ErrNotFound is an response that is used when a resource could not located
// swagger:response errNotFound
func ErrNotFound(errText string, stackTrace error) render.Renderer {
	return &ErrorResponse{
		ErrStack:       stackTrace,
		HTTPStatusCode: http.StatusNotFound,
		ErrorText:      errText,
	}
}

// TODO: will be used if rate limiting will be implemented
// when a user has made too many requests
// ERROR_CAUSE: rate_limiting
func ErrTooManyRequests(err error) render.Renderer {
	return &ErrorResponse{
		ErrStack:       err,
		HTTPStatusCode: http.StatusNotFound,
	}
}

// ERROR CODE 5XX

// when it's not possible to connect to db or
// other dependend external api's
// ERROR_CAUSE: db_connection

// A ErrServiceUnavailable is an response that is used when the API cannot serve any request for reasons such as cannot connect to database etc
// swagger:response errServiceUnavailable
func ErrServiceUnavailable(errText string, stackTrace error) render.Renderer {
	return &ErrorResponse{
		ErrStack:       stackTrace,
		HTTPStatusCode: http.StatusServiceUnavailable,
		ErrorText:      errText,
	}
}

// A ErrInternal is an response that is used when the error that happened doesn't fall in any of the other defined errors
// swagger:response errInternal
func ErrInternal(err error) render.Renderer {
	return &ErrorResponse{
		ErrStack:       err,
		HTTPStatusCode: http.StatusInternalServerError,
		ErrorText:      "Something went wrong",
	}
}
