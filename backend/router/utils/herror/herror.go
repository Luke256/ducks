package herror

import (
	"net/http"

	"github.com/hashicorp/go-multierror"
	"github.com/labstack/echo/v4"
)

func NotFound(message ...any) error {
	return HTTPError(http.StatusNotFound, message)
}

func BadRequest(message ...any) error {
	return HTTPError(http.StatusBadRequest, message)
}

func Forbidden(message ...any) error {
	return HTTPError(http.StatusForbidden, message)
}

func InternalServerError(message ...any) error {
	return HTTPError(http.StatusInternalServerError, message)
}

func Unauthorized(message ...any) error {
	return HTTPError(http.StatusUnauthorized, message)
}

func HTTPError(status int, message any) error {
	switch v := message.(type) {
	case []any:
		if len(v) > 0 {
			return HTTPError(status, v[0])
		}
		return HTTPError(status, nil)

	case string:
		return echo.NewHTTPError(status, v)
		
	case *multierror.Error:
		var errorsList []map[string]string
		for _, err := range v.WrappedErrors() {
			errorsList = append(errorsList, map[string]string{
				"message": err.Error(),
			})
		}
		return echo.NewHTTPError(status, errorsList)

	default:
		return echo.NewHTTPError(status, v)
	}

}