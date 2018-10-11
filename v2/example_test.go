package httpsimp_test

import (
	"log"
	"net/http"
	"net/url"

	"github.com/andreyvit/httpsimplified/v2"
)

const (
	endpointURL = "http://www.example.com/api/v1"
)

func Example() {
	var resp exampleResponse
	// url.Values is just a map[string][]string
	err := httpsimp.Do(httpsimp.MakeGet(endpointURL, "examples/foo.json", url.Values{
		"param1": []string{"value1"},
		"param2": []string{"value2"},
	}, nil), http.DefaultClient, httpsimp.JSON(&resp))

	if err != nil {
		log.Fatal(err)
	}
	log.Printf("foo = %#v", resp)
}

func Example_customHeaders() {
	var resp exampleResponse
	// url.Values and http.Header are both just map[string][]string
	err := httpsimp.Do(httpsimp.MakeGet(endpointURL, "examples/foo.json", url.Values{
		"param1": []string{"value1"},
		"param2": []string{"value2"},
	}, http.Header{
		"X-Powered-By":               []string{"Golang"},
		httpsimp.AuthorizationHeader: []string{httpsimp.BasicAuthValue("user", "secret")},
	}), http.DefaultClient, httpsimp.JSON(&resp))

	if err != nil {
		log.Fatal(err)
	}
	log.Printf("foo = %#v", resp)
}

type exampleResponse struct {
	X string
	Y string
}
