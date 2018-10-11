package httpsimp

import (
	"net/http"
	"net/url"
)

/*
MakeGet builds a GET request with the given URL, headers and params
(encoded into a query string).

base and path are concatenated to form a URL; at least one of them must be
provided, but the other one can be an empty string. The resulting URL must be
valid and parsable via net/url, otherwise panic ensues.

url.Values and http.Header are just maps that can be provided in place,
no need to use their fancy Set or Add methods.
*/
func MakeGet(base, path string, params url.Values, headers http.Header) *http.Request {
	return &http.Request{
		Method: http.MethodGet,
		URL:    URL(base, path, params),
		Header: headers,
	}
}

/*
MakeForm builds a POST/PUT/etc request with the given URL, headers and body
(which contains the given params in application/x-www-form-urlencoded format).

base and path are concatenated to form a URL; at least one of them must be
provided, but the other one can be an empty string. The resulting URL must be
valid and parsable via net/url, otherwise panic ensues.

url.Values and http.Header are just maps that can be provided in place,
no need to use their fancy Set or Add methods.
*/
func MakeForm(method string, base, path string, params url.Values, headers http.Header) *http.Request {
	return EncodeForm(&http.Request{
		Method: method,
		URL:    URL(base, path, nil),
		Header: headers,
	}, params)
}

/*
MakeJSON builds a POST/PUT/etc request with the given URL, headers and body
(which contains the given object encoded in JSON format).

base and path are concatenated to form a URL; at least one of them must be
provided, but the other one can be an empty string. The resulting URL must be
valid and parsable via net/url, otherwise panic ensues.

url.Values and http.Header are just maps that can be provided in place,
no need to use their fancy Set or Add methods.

If JSON encoding fails, the method panics.
*/
func MakeJSON(method string, base, path string, params url.Values, obj interface{}, headers http.Header) *http.Request {
	return EncodeJSONBody(&http.Request{
		Method: method,
		URL:    URL(base, path, params),
		Header: headers,
	}, obj)
}

/*
Make builds a POST/PUT/etc request with the given URL, headers and body.

base and path are concatenated to form a URL; at least one of them must be
provided, but the other one can be an empty string. The resulting URL must be
valid and parsable via net/url, otherwise panic ensues.

url.Values and http.Header are just maps that can be provided in place,
no need to use their fancy Set or Add methods.
*/
func Make(method string, base, path string, params url.Values, body []byte, headers http.Header) *http.Request {
	return SetBody(&http.Request{
		Method: method,
		URL:    URL(base, path, params),
		Header: headers,
	}, body)
}
