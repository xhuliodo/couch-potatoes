package http

import (
	"strings"

	"github.com/go-chi/render"
	"github.com/pkg/errors"
)

const (
	not_authenticated string = "not_authenticated"
	forbidden         string = "forbidden"
	db_connection     string = "db_connection"
	not_found         string = "not_found"
	bad_request       string = "bad_request"
	rate_limiting     string = "rate_limiting"
	app_logic         string = "app_logic"
)

func DetermineErr(err error) render.Renderer {
	cause := errors.Cause(err)

	switch cause.Error() {
	case db_connection:
		textErr := "Service cannot be accessed right now, please try again later"
		return ErrServiceUnavailable(textErr, err)
	case not_found:
		lastErr := getLastErr(err)
		return ErrNotFound(lastErr, err)
	case bad_request:
		lastErr := getLastErr(err)
		return ErrBadRequest(lastErr, err)
	case not_authenticated:
		lastErr := getLastErr(err)
		return ErrUnauthorized(lastErr, err)
	case forbidden:
		lastErr := getLastErr(err)
		return ErrForbidden(lastErr, err)
	default:
		return ErrInternal(err)

	}
}

func getLastErr(err error) string {
	stringErr := err.Error()
	stopChar := ":"
	stopCharIndex := strings.Index(stringErr, stopChar)
	lastErr := stringErr[:stopCharIndex]
	return lastErr
}
