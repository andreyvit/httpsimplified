package httpsimp

import (
	"net/http"
)

type StatusSpec int

const (
	// StatusNone matches no status code and is a zero value of StatusSpec.
	StatusNone StatusSpec = 0

	// StatusAny matches all status codes.
	StatusAny StatusSpec = -1500

	// Status1xx matches all 1xx status codes.
	Status1xx StatusSpec = -100

	// Status2xx matches all 2xx status codes.
	Status2xx StatusSpec = -200

	// Status3xx matches all 3xx status codes.
	Status3xx StatusSpec = -300

	// Status4xx matches all 4xx status codes.
	Status4xx StatusSpec = -400

	// Status4xx matches all 4xx status codes.
	Status5xx StatusSpec = -500

	// Status4xx5xx matches all 4xx and 5xx status codes.
	Status4xx5xx StatusSpec = -900

	StatusOK             = StatusSpec(http.StatusOK)
	StatusCreated        = StatusSpec(http.StatusCreated)
	StatusAccepted       = StatusSpec(http.StatusAccepted)
	StatusNoContent      = StatusSpec(http.StatusNoContent)
	StatusPartialContent = StatusSpec(http.StatusPartialContent)

	StatusUnauthorized = StatusSpec(http.StatusUnauthorized)
	StatusForbidden    = StatusSpec(http.StatusForbidden)
	StatusNotFound     = StatusSpec(http.StatusNotFound)
)

/*
Matches returns whether the given actual HTTP status code matches
the desired status code spec, which may be a specific status code or one
of special constants: StatusNone (won't match anything), Status1xx, Status2xx,
Status3xx, Status4xx, Status5xx.
*/
func (desired StatusSpec) Matches(actual int) bool {
	if actual < 100 || actual > 599 {
		panic("invalid actual status code")
	}

	switch desired {
	case StatusNone:
		return false
	case StatusAny:
		return true
	case Status1xx:
		return (actual >= 100 && actual <= 199)
	case Status2xx:
		return (actual >= 200 && actual <= 299)
	case Status3xx:
		return (actual >= 300 && actual <= 399)
	case Status4xx:
		return (actual >= 400 && actual <= 499)
	case Status5xx:
		return (actual >= 500 && actual <= 599)
	case Status4xx5xx:
		return (actual >= 400 && actual <= 599)
	default:
		if desired < 100 || desired > 599 {
			panic("invalid desired status code spec")
		}
		return actual == int(desired)
	}
}
