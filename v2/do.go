package httpsimp

import (
	"net/http"
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
Do executes the given request via the given http.Client and handles
the body using the specified parsers.

Pass an instance of *http.Client as client. You can use http.DefaultClient,
but note that the default client has no timeouts and might potentially hang
forever, causing goroutine leaks. A custom client is strongly recommended.

For the parsers, use JSON, Bytes, PlainText, Raw or None from this package,
or define your own custom one using MakeParser.
*/
func Do(r *http.Request, client HTTPClient, parsers ...Parser) error {
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
