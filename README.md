# httpsimplified

Package httpsimplified sends outgoing HTTP requests via a simple straightforward API distilled from many internal Golang projects at USA Today Network. It embraces Go stdlib types like url.Values and http.Header, provides composable building blocks for more complex use cases and doesn't try to be clever.

## Install v2

Use v2 for new integrations. It takes an explicit HTTP client so you can set a timeout.

```sh
go get github.com/andreyvit/httpsimplified/v2@latest
```

## Example

This runs against a local test server:

```go
package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	httpsimp "github.com/andreyvit/httpsimplified/v2"
)

func main() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"answer":42}`)
	}))
	defer server.Close()
	client := &http.Client{Timeout: 10 * time.Second}
	var result struct {
		Answer int `json:"answer"`
	}
	err := httpsimp.Do(httpsimp.MakeGet(server.URL, "/", nil, nil), client, httpsimp.JSON(&result))
	if err != nil {
		panic(err)
	}
	fmt.Println(result.Answer) // 42
}
```

See the [v2 guide](v2/) and [API reference](https://pkg.go.dev/github.com/andreyvit/httpsimplified/v2) for request builders, response parsers and authentication.

## Versions

[v1 is deprecated](README-v1.md). Its import path omits `/v2`, and its signatures differ. See the [changelog](CHANGELOG.md) when migrating.

## License

See [LICENSE](LICENSE).
