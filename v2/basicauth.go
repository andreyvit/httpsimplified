package httpsimp

import (
	"encoding/base64"
)

const (
	// AuthorizationHeader is the "Authorization" HTTP header
	AuthorizationHeader = "Authorization"
)

/*
BasicAuthValue returns an Authorization header value for HTTP Basic authentication
method with the given username and password, i.e. it returns:

    "Basic " + base64(username + ":" + password)

Use AuthorizationHeader constant for the header name.
*/
func BasicAuthValue(username, password string) string {
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(username+":"+password))
}
