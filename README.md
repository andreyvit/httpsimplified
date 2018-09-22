# httpsimplified

Package httpsimplified sends outgoing HTTP requests via a simple straightforward API distilled from many internal Golang projects at USA Today Network. It embraces Go stdlib types like url.Values and http.Header, provides composable building blocks for more complex use cases and doesn't try to be clever.

```go
client := &http.Client{
    Timeout: time.Second * 10,
}

var resp responseType
err := httpsimplified.Get(baseURL, path, params, headers, client, httpsimplified.JSON(&resp))
```

## v2

[![README v2](https://img.shields.io/badge/readme-v2-green.svg)](v2/) [![GoDoc](https://godoc.org/github.com/andreyvit/httpsimplified/v2?status.svg)](https://godoc.org/github.com/andreyvit/httpsimplified/v2)

I've learned that [using `http.DefaultClient` is a bad idea](https://medium.com/@nate510/don-t-use-go-s-default-http-client-4804cb19f779):

> TL;DR: Go’s http package doesn’t specify request timeouts by default, allowing services to hijack your goroutines. Always specify a custom http.Client when connecting to outside services.

So all functions in v2 accept an additional `*http.Client` argument. You can pass `http.DefaultClient` if you want the old behavior. See [CHANGELOG](CHANGELOG.md) for all changes.


## v1

[![README v2](https://img.shields.io/badge/readme-v1-orange.svg)](README-v1.md) [![GoDoc](https://godoc.org/github.com/andreyvit/httpsimplified?status.svg)](https://godoc.org/github.com/andreyvit/httpsimplified)

This version is deprecated.



See [godoc.org/github.com/andreyvit/httpsimplified](https://godoc.org/github.com/andreyvit/httpsimplified) for a full reference.

Call Get, Post or Put to send a request and parse the response in a single call:


where httpsimplified.JSON is a body parser function (we also provide PlainText, Bytes, Raw and None parsers, and you can define your own). See the example for more details.

For more advanced requests, build http.Request yourself and call Perform:

```go
var resp responseType
err := httpsimplified.Perform(&http.Request{
    Method: http.MethodPut,
    URL:    httpsimplified.URL(baseURL, path, params),
    Header: http.Header{...},
    Body:   []byte{"whatever"},
}, httpsimplified.JSON, &resp)
```

Use URL func to concatenate a URL and include query params.

Use EncodeBody helper to generate application/x-www-form-urlencoded bodies.

Finally, if http.DefaultClient doesn't rock your boat, you're free to build and execute a request through whatever means necessary and then call JSON, Bytes or None to verify the response status code and handle the body:

```go
req := EncodeBody(&http.Request{
    Method: http.MethodPost,
    URL:    httpsimplified.URL(baseURL, path, nil),
    Header: http.Header{...},
}, url.Params{...})

httpResp, err := myCustomClient.Do(req)
if err != nil { ... }

var resp responseType
err = httpsimplified.JSON(httpResp, &resp)
```

To handle HTTP basic authentication, use BasicAuth helper:

```go
err := httpsimplified.Get("...", "...", url.Values{...}, http.Header{
    httpsimplified.AuthorizationHeader: []string{httpsimplified.BasicAuth("user", "pw")},
}, httpsimplified.JSON, &resp)
```
