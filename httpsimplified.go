/*
Package httpsimplified sends outgoing HTTP requests via a simple straightforward
API distilled from many internal Golang projects at USA Today Network.
It embraces Go stdlib types like url.Values and http.Header, provides composable
building blocks for more complex use cases and doesn't try to be clever.

Call Get, Post or Put to send a request and parse the response in a single call:

	var resp responseType
	err := httpsimplified.Get(baseURL, path, params, headers, httpsimplified.JSON, &resp)

where httpsimplified.JSON is a body parser function (we also provide Bytes, Raw
and None parsers, and you can define your own). See the example for more details.

For more advanced requests, build http.Request yourself and call Perform:

	var resp responseType
	err := httpsimplified.Perform(&http.Request{
		Method: http.MethodPut,
		URL:    httpsimplified.URL(baseURL, path, params),
		Header: http.Header{...},
		Body:   []byte{"whatever"},
	}, httpsimplified.JSON, &resp)

Use URL func to concatenate a URL and include query params.

Use EncodeBody helper to generate application/x-www-form-urlencoded bodies.

Finally, if http.DefaultClient doesn't rock your boat, you're free to build and
execute a request through whatever means necessary and then call JSON, Bytes or
None to verify the response status code and handle the body:

	req := EncodeBody(&http.Request{
		Method: http.MethodPost,
		URL:    httpsimplified.URL(baseURL, path, nil),
		Header: http.Header{...},
	}, url.Params{...})

	httpResp, err := myCustomClient.Do(req)
	if err != nil { ... }

	var resp responseType
	err = httpsimplified.JSON(httpResp, &resp)

To handle HTTP basic authentication, use BasicAuth helper:

	err := httpsimplified.Get("...", "...", url.Values{...}, http.Header{
		httpsimplified.AuthorizationHeader: []string{httpsimplified.BasicAuth("user", "pw")},
	}, httpsimplified.JSON, &resp)
*/
package httpsimplified

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"net/url"
	"strings"
)

const (
	// ContentTypeJSON is "application/json"
	ContentTypeJSON = "application/json"

	// ContentTypeFormURLEncoded is "application/x-www-form-urlencoded"
	ContentTypeFormURLEncoded = "application/x-www-form-urlencoded"

	// AuthorizationHeader is the "Authorization" HTTP header
	AuthorizationHeader = "Authorization"
)

type Error struct {
	Method string
	Path   string
	Cause  error
}

func (err *Error) Error() string {
	return fmt.Sprintf("%s %s: %v", err.Method, err.Path, err.Cause)
}

type StatusError struct {
	StatusCode int

	ContentType string

	Body interface{}

	DecodingError error
}

func (err *StatusError) Error() string {
	if err.ContentType == ContentTypeJSON {
		if err.DecodingError != nil {
			return fmt.Sprintf("HTTP %d, error decoding JSON: %v", err.StatusCode, err.DecodingError)
		} else {
			return fmt.Sprintf("HTTP %d, JSON: %v", err.StatusCode, err.Body)
		}
	} else {
		return fmt.Sprintf("HTTP %d, response type: %v", err.StatusCode, err.ContentType)
	}
}

type ContentTypeError struct {
	StatusCode int

	ContentType string

	ExpectedContentType string
}

func (err *ContentTypeError) Error() string {
	return fmt.Sprintf("HTTP %s, but unexpected response type %v, wanted %v", err.StatusCode, err.ContentType, err.ExpectedContentType)
}

func CheckStatusError(err error) *StatusError {
	if e, ok := err.(*Error); ok {
		err = e.Cause
	}

	e, _ := err.(*StatusError)
	return e
}

func verify(resp *http.Response, expectedCType string) error {
	mediaType := resp.Header.Get("Content-Type")
	ctype, _, err := mime.ParseMediaType(mediaType)
	if err != nil {
		return fmt.Errorf("cannot parse Content-Type string %+v", mediaType)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var body interface{}
		if ctype == ContentTypeJSON {
			err = json.NewDecoder(resp.Body).Decode(&body)
			if err == nil {
				return &StatusError{resp.StatusCode, ctype, body, nil}
			} else {
				return &StatusError{resp.StatusCode, ctype, nil, err}
			}
		} else {
			return &StatusError{resp.StatusCode, ctype, nil, nil}
		}
	}

	if expectedCType != "" && ctype != expectedCType {
		return &ContentTypeError{resp.StatusCode, ctype, expectedCType}
	}

	return nil
}

/*
Parser is a function used to verify and handle the HTTP response. This package
provides a number of parser functions, and you can define your own.
*/
type Parser func(resp *http.Response, result interface{}) error

/*
Raw is a Parser function that verifies the response status code and returns
the raw *http.Response; result must be a pointer to *http.Response variable.
*/
func Raw(resp *http.Response, result interface{}) error {
	err := verify(resp, "")
	if err != nil {
		resp.Body.Close()
		return err
	}

	ptr := result.(**http.Response)
	*ptr = resp
	return nil
}

/*
JSON is a Parser function that verifies the response status code and content
type (which must be ContentTypeJSON) and unmarshals the body into the
result variable (which can be anything that you'd pass to json.Unmarshal).
*/
func JSON(resp *http.Response, result interface{}) error {
	defer resp.Body.Close()

	err := verify(resp, ContentTypeJSON)
	if err != nil {
		return err
	}

	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		return fmt.Errorf("HTTP %s, error decoding JSON: %v", resp.Status, err)
	}

	return nil
}

/*
Bytes is a Parser function that verifies the response status code and reads
the entire body into a byte array; result must be a pointer to a []byte variable.
*/
func Bytes(resp *http.Response, result interface{}) error {
	defer resp.Body.Close()

	err := verify(resp, "")
	if err != nil {
		return err
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("HTTP %s, error reading body: %v", resp.Status, err)
	}

	*(result.(*[]byte)) = b
	return nil
}

/*
None is a Parser function that verifies the response status code and discards
the response body; result argument is ignored and should be nil.

A typical use would be to pass this function into Get, Post or Perform,
but you can also call it directly.
*/
func None(resp *http.Response, result interface{}) error {
	resp.Body.Close()
	return verify(resp, "")
}

/*
Get builds a GET request with the given URL, parameters and headers, executes
it via http.DefaultClient.Do and handles the body using the specified parser
function.

base and path are concatenated to form a URL; at least one of them must be
provided, but the other one can be an empty string. The resulting URL must be
valid and parsable via net/url, otherwise panic ensues.

url.Values and http.Header are just maps that can be provided in place,
no need to use their fancy Set or Add methods.

parser can be either JSON, Bytes, Raw or None from this package,
or your own custom parser function; it will be called with *http.Response and
the result you pass in.
*/
func Get(base, path string, params url.Values, headers http.Header, parser Parser, result interface{}) error {
	return Perform(&http.Request{
		Method: http.MethodGet,
		URL:    URL(base, path, params),
		Header: headers,
	}, parser, result)
}

/*
Post builds a POST request with the given URL, headers and body (which contains
the given params in application/x-www-form-urlencoded format), executes
it via http.DefaultClient.Do and handles the body using the specified parser
function.

base and path are concatenated to form a URL; at least one of them must be
provided, but the other one can be an empty string. The resulting URL must be
valid and parsable via net/url, otherwise panic ensues.

url.Values and http.Header are just maps that can be provided in place,
no need to use their fancy Set or Add methods.

parser can be either JSON, Bytes, Raw or None from this package,
or your own custom parser function; it will be called with *http.Response and
the result you pass in.
*/
func Post(base, path string, params url.Values, headers http.Header, parser Parser, result interface{}) error {
	return Perform(EncodeBody(&http.Request{
		Method: http.MethodPost,
		URL:    URL(base, path, nil),
		Header: headers,
	}, params), parser, result)
}

/*
Put builds a PUT request with the given URL, headers and body (which contains
the given params in application/x-www-form-urlencoded format), executes
it via http.DefaultClient.Do and handles the body using the specified parser
function.

base and path are concatenated to form a URL; at least one of them must be
provided, but the other one can be an empty string. The resulting URL must be
valid and parsable via net/url, otherwise panic ensues.

url.Values and http.Header are just maps that can be provided in place,
no need to use their fancy Set or Add methods.

parser can be either JSON, Bytes, Raw or None from this package,
or your own custom parser function; it will be called with *http.Response and
the result you pass in.
*/
func Put(base, path string, params url.Values, headers http.Header, parser Parser, result interface{}) error {
	return Perform(EncodeBody(&http.Request{
		Method: http.MethodPut,
		URL:    URL(base, path, nil),
		Header: headers,
	}, params), parser, result)
}

/*
Perform executes the given request via http.DefaultClient.Do and handles
the body using the specified parser function.

parser can be either JSON, Bytes, Raw or None from this package,
or your own custom parser function; it will be called with *http.Response and
the result you pass in.
*/
func Perform(r *http.Request, parser Parser, result interface{}) error {
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return &Error{r.Method, r.URL.Path, err}
	}

	err = parser(resp, result)
	if err != nil {
		return &Error{r.Method, r.URL.Path, err}
	}

	return nil
}

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
EncodeBody encodes the given params into application/x-www-form-urlencoded
format and sets the body and Content-Type on the given request.

To properly handle HTTP redirects, both Body and GetBody are set.
*/
func EncodeBody(r *http.Request, params url.Values) *http.Request {
	if params == nil {
		params = url.Values{}
	}
	body := []byte(params.Encode())

	r.Body = ioutil.NopCloser(bytes.NewReader(body))
	r.GetBody = func() (io.ReadCloser, error) {
		return ioutil.NopCloser(bytes.NewReader(body)), nil
	}

	if r.Header == nil {
		r.Header = make(http.Header)
	}
	if r.Header["Content-Type"] == nil {
		r.Header["Content-Type"] = []string{ContentTypeFormURLEncoded}
	}

	return r
}

/*
BasicAuth returns an Authorization header value for HTTP Basic authentication
method with the given username and password, i.e. it returns:

	"Basic " + base64(username + ":" + password)

Use AuthorizationHeader constant for the header name.
*/
func BasicAuth(username, password string) string {
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(username+":"+password))
}
