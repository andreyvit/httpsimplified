package httpsimp

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"unicode/utf8"
)

/*
Raw is a Parser function that verifies the response status code and returns
the raw *http.Response; result must be a pointer to *http.Response variable.
*/
func Raw(ptr **http.Response, mopt ...ParseOption) Parser {
	return MakeParser("", mopt, func(resp *http.Response) (interface{}, error) {
		*ptr = resp
		return nil, nil
	})
}

/*
JSON is a Parser function that verifies the response status code and content
type (which must be ContentTypeJSON) and unmarshals the body into the
result variable (which can be anything that you'd pass to json.Unmarshal).
*/
func JSON(result interface{}, mopt ...ParseOption) Parser {
	if result == nil {
		var body interface{}
		result = &body
	}
	return MakeParser(ContentTypeJSON, mopt, func(resp *http.Response) (interface{}, error) {
		defer resp.Body.Close()
		err := json.NewDecoder(resp.Body).Decode(result)
		body := reflect.ValueOf(result).Elem().Interface()
		return body, err
	})
}

/*
Bytes is a Parser function that verifies the response status code and reads
the entire body into a byte array; result must be a pointer to a []byte variable.
*/
func Bytes(result *[]byte, mopt ...ParseOption) Parser {
	return MakeParser("", mopt, func(resp *http.Response) (interface{}, error) {
		defer resp.Body.Close()
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			err = fmt.Errorf("error reading body: %v", err)
		}
		*result = b
		return b, err
	})
}

/*
PlainText is a Parser function that verifies the response status code and reads
the entire body into a string; result must be a pointer to a string variable.
*/
func PlainText(result *string, mopt ...ParseOption) Parser {
	if result == nil {
		var body string
		result = &body
	}
	return MakeParser("", mopt, func(resp *http.Response) (interface{}, error) {
		defer resp.Body.Close()
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			err = fmt.Errorf("error reading body: %v", err)
		}
		if !utf8.Valid(b) {
			return b, errors.New("invalid utf-8 sequence encountered")
		}

		s := string(b)
		*result = s
		return s, err
	})
}

/*
None is a Parser function that verifies the response status code and discards
the response body; result argument is ignored and should be nil.

A typical use would be to pass this function into Get, Post or Perform,
but you can also call it directly.
*/
func None(mopt ...ParseOption) Parser {
	return MakeParser("", mopt, func(resp *http.Response) (interface{}, error) {
		resp.Body.Close()
		return nil, nil
	})
}
