# httpsimplified

[![GoDoc](https://godoc.org/github.com/andreyvit/httpsimplified/v2?status.svg)](https://godoc.org/github.com/andreyvit/httpsimplified/v2)

Package httpsimplified sends outgoing HTTP requests via a simple straightforward API distilled from many internal Golang projects at USA Today Network. It embraces Go stdlib types like url.Values and http.Header, provides composable building blocks for more complex use cases and doesn't try to be clever.

See [godoc.org/github.com/andreyvit/httpsimplified/v2](https://godoc.org/github.com/andreyvit/httpsimplified/v2) for a full reference.

Call Get, Post or Put to send a request and parse the response in a single call:

```go
var resp responseType
err := httpsimplified.Get(baseURL, path, params, headers, clients, httpsimplified.JSON(&resp))
```

where httpsimplified.JSON is a body parser function (we also provide PlainText, Bytes, Raw and None parsers, and you can define your own). See the examples and GoDoc for more details.

See [CHANGELOG](../CHANGELOG.md) for a history of changes.
