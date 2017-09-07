package httpsimplified

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"net/url"
	"strings"
)

const (
	ContentTypeJSON           = "application/json"
	ContentTypeFormURLEncoded = "application/x-www-form-urlencoded"
)

func verify(resp *http.Response, expectedCType string) error {
	mediaType := resp.Header.Get("Content-Type")
	ctype, _, err := mime.ParseMediaType(mediaType)
	if err != nil {
		return fmt.Errorf("cannot parse Content-Type string %+v", mediaType)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var body interface{}
		if ctype == ContentTypeJSON {
			err = json.NewDecoder(resp.Body).Decode(&body)
			if err == nil {
				return fmt.Errorf("HTTP Status %s, JSON %v", resp.Status, body)
			} else {
				return fmt.Errorf("HTTP Status %s, error decoding JSON: %v", resp.Status, err)
			}
		} else {
			return fmt.Errorf("HTTP Status %s, %v response", resp.Status, ctype)
		}
	}

	if expectedCType != "" && ctype != expectedCType {
		return fmt.Errorf("HTTP Status %s, but unexpected response type %v, wanted %v", resp.Status, ctype, expectedCType)
	}

	return nil
}

type Parser func(resp *http.Response, result interface{}) error

func Raw(resp *http.Response, result interface{}) error {
	err := verify(resp, "")
	if err != nil {
		resp.Body.Close()
		return err
	}

	ptr := result.(**http.Response)
	*ptr = resp
	return nil
}

func JSON(resp *http.Response, result interface{}) error {
	defer resp.Body.Close()

	err := verify(resp, ContentTypeJSON)
	if err != nil {
		return err
	}

	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		return fmt.Errorf("HTTP %s, error decoding JSON: %v", resp.Status, err)
	}

	return nil
}

func None(resp *http.Response, result interface{}) error {
	resp.Body.Close()
	return verify(resp, "")
}

func Get(base, path string, params url.Values, headers http.Header, parser Parser, result interface{}) error {
	return Perform(&http.Request{
		Method: http.MethodGet,
		URL:    URL(base, path, params),
		Header: headers,
	}, parser, result)
}

func Post(base, path string, params url.Values, headers http.Header, parser Parser, result interface{}) error {
	return Perform(EncodeBody(&http.Request{
		Method: http.MethodPost,
		URL:    URL(base, path, nil),
		Header: headers,
	}, params), parser, result)
}

func Put(base, path string, params url.Values, headers http.Header, parser Parser, result interface{}) error {
	return Perform(EncodeBody(&http.Request{
		Method: http.MethodPut,
		URL:    URL(base, path, nil),
		Header: headers,
	}, params), parser, result)
}

func Perform(r *http.Request, parser Parser, result interface{}) error {
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return fmt.Errorf("%s %s: %v", r.Method, r.URL.Path, err)
	}

	err = parser(resp, result)
	if err != nil {
		return fmt.Errorf("%s %s: %v", r.Method, r.URL.Path, err)
	}

	return nil
}

func URL(base, path string, params url.Values) *url.URL {
	var components *url.URL
	var err error

	if base == "" {
		components, err = url.Parse(path)
		if err != nil {
			panic(err)
		}
	} else {
		components, err = url.Parse(base)
		if err != nil {
			panic(err)
		}

		if path != "" {
			if !strings.HasPrefix(path, "/") {
				path = "/" + path
			}
			components.Path = components.Path + path
		}
	}

	if params != nil {
		components.RawQuery = params.Encode()
	}

	return components
}

func EncodeBody(r *http.Request, params url.Values) *http.Request {
	if params == nil {
		params = url.Values{}
	}
	body := []byte(params.Encode())

	r.Body = ioutil.NopCloser(bytes.NewReader(body))
	r.GetBody = func() (io.ReadCloser, error) {
		return ioutil.NopCloser(bytes.NewReader(body)), nil
	}

	if r.Header == nil {
		r.Header = make(http.Header)
	}
	if r.Header["Content-Type"] == nil {
		r.Header["Content-Type"] = []string{ContentTypeFormURLEncoded}
	}

	return r
}
