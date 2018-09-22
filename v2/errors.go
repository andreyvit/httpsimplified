package httpsimp

import (
	"fmt"
)

type wrapperError struct {
	Method string
	Path   string
	Cause  error
}

func (err *wrapperError) Error() string {
	if err.Path != "" {
		return fmt.Sprintf("%s %s: %v", err.Method, err.Path, err.Cause)
	} else {
		return fmt.Sprintf("%s: %v", err.Method, err.Cause)
	}
}

type responseError struct {
	StatusCode int

	ContentType       string
	WantedContentType string
	ContentTypeOK     bool

	Body          interface{}
	DecodingError error
}

func (err *responseError) Error() string {
	if !err.ContentTypeOK {
		if err.DecodingError != nil {
			return fmt.Sprintf("HTTP %d, unexpected response of type %v, wanted %v; error decoding response body: %v", err.StatusCode, err.ContentType, err.WantedContentType, err.DecodingError)
		} else if err.Body != nil {
			return fmt.Sprintf("HTTP %d, unexpected response of type %v, wanted %v: %v", err.StatusCode, err.ContentType, err.WantedContentType, err.Body)
		} else {
			return fmt.Sprintf("HTTP %d, unexpected response type %v, wanted %v", err.StatusCode, err.ContentType, err.WantedContentType)
		}
	} else {
		if err.DecodingError != nil {
			return fmt.Sprintf("HTTP %d, error decoding %v response: %v", err.StatusCode, err.ContentType, err.DecodingError)
		} else if err.Body != nil {
			return fmt.Sprintf("HTTP %d, %v response: %v", err.StatusCode, err.ContentType, err.Body)
		} else {
			return fmt.Sprintf("HTTP %d, %v response", err.StatusCode, err.ContentType)
		}
	}
}

func getResponseError(err error) *responseError {
	if e, ok := err.(*wrapperError); ok {
		err = e.Cause
	}

	e, _ := err.(*responseError)
	return e
}

func StatusCode(err error) int {
	if e := getResponseError(err); e != nil {
		return e.StatusCode
	} else {
		return 0
	}
}

func Is5xx(err error) bool {
	code := StatusCode(err)
	return (code != 0) && (code >= 500 && code <= 599)
}

func Is4xx(err error) bool {
	code := StatusCode(err)
	return (code != 0) && (code >= 400 && code <= 499)
}
