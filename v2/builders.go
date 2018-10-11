package httpsimp

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

/*
URL returns a *url.URL (conveniently suitable for http.Request's URL field)
concatenating the two given URL strings and optionally appending a query string
with the given params.

base and path are concatenated to form a URL; at least one of them must be
provided, but the other one can be an empty string. The resulting URL must be
valid and parsable via net/url, otherwise panic ensues.

url.Values and http.Header are just maps that can be provided in place,
no need to use their fancy Set or Add methods.
*/
func URL(base, path string, params url.Values) *url.URL {
	var components *url.URL
	var err error

	if base == "" {
		components, err = url.Parse(path)
		if err != nil {
			panic(err)
		}
	} else {
		components, err = url.Parse(base)
		if err != nil {
			panic(err)
		}

		if path != "" {
			if !strings.HasPrefix(path, "/") {
				path = "/" + path
			}
			components.Path = components.Path + path
		}
	}

	if params != nil {
		components.RawQuery = strings.Replace(params.Encode(), "+", "%20", -1)
	}

	return components
}

/*
EncodeForm encodes the given params into application/x-www-form-urlencoded
format and sets the body and Content-Type on the given request.

To properly handle HTTP redirects, both Body and GetBody are set.
*/
func EncodeForm(r *http.Request, params url.Values) *http.Request {
	if params == nil {
		params = url.Values{}
	}
	body := []byte(params.Encode())
	_ = SetBody(r, body)

	if r.Header == nil {
		r.Header = make(http.Header)
	}
	if r.Header["Content-Type"] == nil {
		r.Header["Content-Type"] = []string{ContentTypeFormURLEncoded}
	}

	return r
}

/*
EncodeJSONBody encodes the given object into JSON (application/json)
format and sets the body and Content-Type on the given request.

If JSON encoding fails, the method panics.

To properly handle HTTP redirects, both Body and GetBody are set.
*/
func EncodeJSONBody(r *http.Request, obj interface{}) *http.Request {
	body, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	_ = SetBody(r, body)

	if r.Header == nil {
		r.Header = make(http.Header)
	}
	if r.Header["Content-Type"] == nil {
		r.Header["Content-Type"] = []string{ContentTypeJSON}
	}

	return r
}

/*
SetBody sets the given request's body to the given data.

To properly handle HTTP redirects, both Body and GetBody are set.
*/
func SetBody(r *http.Request, data []byte) *http.Request {
	r.Body = ioutil.NopCloser(bytes.NewReader(data))
	r.GetBody = func() (io.ReadCloser, error) {
		return ioutil.NopCloser(bytes.NewReader(data)), nil
	}
	return r
}
