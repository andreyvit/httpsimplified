package httpsimp

import (
	"net/http"
	"net/url"
)

/*
HTTPClient is an interface implemented by *http.Client, requiring
only the Do method. Instead of accepting *http.Client, the methods
in this package accept HTTPClients for extra flexibility.
*/
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

/*
Get builds a GET request with the given URL, parameters and headers, executes
it via the given http.Client.Do and handles the body using the specified parser
functions.

base and path are concatenated to form a URL; at least one of them must be
provided, but the other one can be an empty string. The resulting URL must be
valid and parsable via net/url, otherwise panic ensues.

url.Values and http.Header are just maps that can be provided in place,
no need to use their fancy Set or Add methods.

Pass an instance of *http.Client as client. You can use http.DefaultClient,
but note that the default client has no timeouts and might potentially hang
forever, causing goroutine leaks. A custom client is strongly recommended.

For the parsers, use JSON, Bytes, PlainText, Raw or None from this package,
or define your own custom one using MakeParser.
*/
func Get(base, path string, params url.Values, headers http.Header, client HTTPClient, parsers ...Parser) error {
	return Perform(&http.Request{
		Method: http.MethodGet,
		URL:    URL(base, path, params),
		Header: headers,
	}, client, parsers...)
}

/*
Post builds a POST request with the given URL, headers and body (which contains
the given params in application/x-www-form-urlencoded format), executes
it via the given http.Client.Do and handles the body using the specified parsers.

base and path are concatenated to form a URL; at least one of them must be
provided, but the other one can be an empty string. The resulting URL must be
valid and parsable via net/url, otherwise panic ensues.

url.Values and http.Header are just maps that can be provided in place,
no need to use their fancy Set or Add methods.

Pass an instance of *http.Client as client. You can use http.DefaultClient,
but note that the default client has no timeouts and might potentially hang
forever, causing goroutine leaks. A custom client is strongly recommended.

For the parsers, use JSON, Bytes, PlainText, Raw or None from this package,
or define your own custom one using MakeParser.
*/
func Post(base, path string, params url.Values, headers http.Header, client HTTPClient, parsers ...Parser) error {
	return Perform(EncodeForm(&http.Request{
		Method: http.MethodPost,
		URL:    URL(base, path, nil),
		Header: headers,
	}, params), client, parsers...)
}

/*
Put builds a PUT request with the given URL, headers and body (which contains
the given params in application/x-www-form-urlencoded format), executes
it via the given http.Client and handles the body using the specified parsers.

base and path are concatenated to form a URL; at least one of them must be
provided, but the other one can be an empty string. The resulting URL must be
valid and parsable via net/url, otherwise panic ensues.

url.Values and http.Header are just maps that can be provided in place,
no need to use their fancy Set or Add methods.

Pass an instance of *http.Client as client. You can use http.DefaultClient,
but note that the default client has no timeouts and might potentially hang
forever, causing goroutine leaks. A custom client is strongly recommended.

For the parsers, use JSON, Bytes, PlainText, Raw or None from this package,
or define your own custom one using MakeParser.
*/
func Put(base, path string, params url.Values, headers http.Header, client HTTPClient, parsers ...Parser) error {
	return Perform(EncodeForm(&http.Request{
		Method: http.MethodPut,
		URL:    URL(base, path, nil),
		Header: headers,
	}, params), client, parsers...)
}

/*
Perform executes the given request via the given http.Client and handles
the body using the specified parsers.

Pass an instance of *http.Client as client. You can use http.DefaultClient,
but note that the default client has no timeouts and might potentially hang
forever, causing goroutine leaks. A custom client is strongly recommended.

For the parsers, use JSON, Bytes, PlainText, Raw or None from this package,
or define your own custom one using MakeParser.
*/
func Perform(r *http.Request, client HTTPClient, parsers ...Parser) error {
	resp, err := client.Do(r)
	if err != nil {
		return &wrapperError{r.Method, r.URL.Path, err}
	}

	err = Parse(resp, parsers...)
	if err != nil {
		return &wrapperError{r.Method, r.URL.Path, err}
	}

	return nil
}
