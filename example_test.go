package httpsimplified_test

import (
	"log"
	"net/http"
	"net/url"

	"github.com/andreyvit/httpsimplified"
)

const (
	endpointURL = "http://www.example.com/api/v1"
)

func Example() {
	var resp exampleResponse
	// url.Values is just a map[string][]string
	err := httpsimplified.Get(endpointURL, "examples/foo.json", url.Values{
		"param1": []string{"value1"},
		"param2": []string{"value2"},
	}, nil, httpsimplified.JSON, &resp)

	if err != nil {
		log.Fatal(err)
	}
	log.Printf("foo = %#v", resp)
}

func Example_customHeaders() {
	var resp exampleResponse
	// url.Values and http.Header are both just map[string][]string
	err := httpsimplified.Get(endpointURL, "examples/foo.json", url.Values{
		"param1": []string{"value1"},
		"param2": []string{"value2"},
	}, http.Header{
		"X-Powered-By":                     []string{"Golang"},
		httpsimplified.AuthorizationHeader: []string{httpsimplified.BasicAuth("user", "secret")},
	}, httpsimplified.JSON, &resp)

	if err != nil {
		log.Fatal(err)
	}
	log.Printf("foo = %#v", resp)
}

type exampleResponse struct {
	X string
	Y string
}
