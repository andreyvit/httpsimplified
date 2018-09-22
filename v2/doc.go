/*
Package httpsimp sends outgoing HTTP requests via a simple straightforward
API distilled from many internal Golang projects at USA Today Network.
It embraces Go stdlib types like url.Values and http.Header, provides composable
building blocks for more complex use cases and doesn't try to be clever.

Call Get, Post or Put to send a request and parse the response in a single call:

    var resp responseType
    err := httpsimp.Get(baseURL, path, params, headers, client, httpsimp.JSON(&resp))

where httpsimp.JSON is a body parser function (we also provide PlainText,
Bytes, Raw and None parsers, and you can define your own).
See the example for more details.

You can adjust body parser parameters by passing additional options to body
parser functions, like this:

    httpsimp.JSON(nil, httpsimp.ContentType("application/something"))

Available options:

- httpsimp.StatusAny, httpsimp.Status4xx, httpsimp.Status4xx5xx, or a specific
status like httpsimp.StatusOK or httpsimp.StatusSpec(http.StatusTeapot) will
match only responses with the given status.

- httpsimp.ContentType("application/something") will match only response with
the given content type.

- httpsimp.ContentType("") will match any content type (can be used to cancel
default application/json filter used by JSON).

- httpsimp.ReturnError() results in a non-nil error returned.

Pass multiple parsers to handle alternative response types or non-2xx status codes:

    var resp responseStruct
    var bytes []byte
    var e errorStruct
    err := httpsimp.Get(...,
        httpsimp.JSON(&resp),
        httpsimp.Bytes(&bytes, httpsimp.ContentType("image/png")),
        httpsimp.JSON(&e, httpsimp.Status4xx5xx))

For more advanced requests, build http.Request yourself and call Perform:

    var resp responseType
    err := httpsimp.Perform(&http.Request{
        Method: http.MethodPut,
        URL:    httpsimp.URL(baseURL, path, params),
        Header: http.Header{...},
        Body:   []byte{"whatever"},
    }, httpsimp.JSON(&resp))

Use URL func to concatenate a URL and include query params, and EncodeForm
helper to generate application/x-www-form-urlencoded bodies.

Finally, if http.Client.Do doesn't rock your boat, you're free to build and
execute a request through whatever means necessary and then call Parse
to verify the response status code and handle the body:

    req := EncodeForm(&http.Request{
        Method: http.MethodPost,
        URL:    httpsimp.URL(baseURL, path, nil),
        Header: http.Header{...},
    }, url.Params{...})

    httpResp, err := whatever.Do(req)
    if err != nil { ... }

    var resp responseType
    err = httpsimp.Parse(httpResp, httpsimp.JSON(&resp))

To handle HTTP basic authentication, use BasicAuthValue helper:

    err := httpsimp.Get("...", "...", url.Values{...}, http.Header{
        httpsimp.AuthorizationHeader: []string{httpsimp.BasicAuthValue("user", "pw")},
    }, httpsimp.JSON, &resp)
*/
package httpsimp
