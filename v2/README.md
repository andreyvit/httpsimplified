# httpsimplified

[![GoDoc](https://godoc.org/github.com/andreyvit/httpsimplified/v2?status.svg)](https://godoc.org/github.com/andreyvit/httpsimplified/v2)

Package httpsimplified sends outgoing HTTP requests via a simple straightforward API distilled from many internal Golang projects at USA Today Network. It embraces Go stdlib types like url.Values and http.Header, provides composable building blocks for more complex use cases and doesn't try to be clever.

See [godoc.org/github.com/andreyvit/httpsimplified/v2](https://godoc.org/github.com/andreyvit/httpsimplified/v2) for a full reference.


Usage
-----

Call `Do` with `MakeGet`, `MakeForm`, `MakeJSON` or `Make` to send a request and parse the response:

```go
var resp responseType
err := httpsimp.Do(httpsimp.MakeGet(baseURL, path, params, headers), client, httpsimp.JSON(&resp))
```

where:

* `httpsimp.MakeGet` is a request builder function returning `*http.Request`, we also provide `MakeForm`, `MakeJSON` and `Make` (and, of course, you're free to build `http.Request` any other way);

* `httpsimp.JSON` is a body parser, we also provide `PlainText`, `Bytes`, `Raw` and `None` parsers; you can use multiple parsers and/or adjust their parameters; and you can define your own parsers.

See [GoDoc](https://godoc.org/github.com/andreyvit/httpsimplified/v2) for more details.


Features
--------

Provides simple, composable building blocks and embraces Go stdlib types:

* request builder functions (`MakeGet`, `MakeForm`, `MakeJSON`, `Make`) just return an `*http.Request` that you can further customize if you want (e.g. you can call `.WithContext(ctx)` to make the request cancelable);

* `Parse` parses any `*http.Response` using one or more body parsers;

* `Do` accepts any `*http.Request`, executes it using `http.Client` and handles the response via `Parse`;

* when building a custom request, you can use lower-level helper functions: `URL`, `EncodeForm`, `EncodeJSONBody`, `SetBody`;

* request parameters are specified via `url.Values` — you can pass an inline map like `url.Values{"key": []string{"value"}}`, or build `url.Values` some other way, or just pass `nil`;

* request headers are specified via `http.Header` — you can pass an inline map like `http.Header{"X-Something": []string{"value"}}`, or build `http.Header` some other way, or just pass `nil`.

This library strives to be as straight-forward and non-magical as possible.

We used to define one-shot helpers like `Get` and `Post`, but then the number of combination exploded (`Get`, `GetContext`, `PostForm`, `PostFormContext`, `PostJSON`, `PostJSONContext`, plus same for PUT), and we instead opted for a 2-call combination (`Do(MakeXxx(http.MethodPost, ...), ...)`), which we believe to be superior in every way.


Change Log
----------

See [CHANGELOG](../CHANGELOG.md) for a history of changes.
