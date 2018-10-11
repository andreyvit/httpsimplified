package httpsimp

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func get(statusCode int, ctype string, body []byte, parsers ...Parser) error {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", ctype)
		w.WriteHeader(statusCode)
		w.Write(body)
	}))
	defer srv.Close()

	return Do(MakeGet("", srv.URL, nil, nil), http.DefaultClient, parsers...)
}

func TestGetJSON200(t *testing.T) {
	var resp struct {
		Foo int `json:"foo"`
	}
	err := get(http.StatusOK, ContentTypeJSON, []byte(`{"foo": 42}`), JSON(&resp))
	if err != nil {
		t.Fatal(err)
	}
	if resp.Foo != 42 {
		t.Fatalf("invalid value of Foo: %v", resp)
	}
}

func TestGetInvalidContentType200(t *testing.T) {
	var resp struct {
		Foo int `json:"foo"`
	}
	err := get(http.StatusOK, ContentTypeTextPlain, []byte(`{"foo": 42}`), JSON(&resp))
	if err == nil {
		t.Fatal("err is nil")
	}
	if !strings.Contains(err.Error(), `HTTP 200, unexpected response type text/plain, wanted application/json`) {
		t.Fatalf("invalid error: %v", err)
	}
}

func TestGetAltJSON200(t *testing.T) {
	var resp struct {
		Foo int `json:"foo"`
	}
	var text string
	err := get(http.StatusOK, ContentTypeJSON, []byte(`{"foo": 42}`), JSON(&resp), PlainText(&text))
	if err != nil {
		t.Fatal(err)
	}
	if resp.Foo != 42 {
		t.Fatalf("invalid value of Foo: %v", resp)
	}
}

func TestGetAltText200(t *testing.T) {
	var text string
	err := get(http.StatusOK, ContentTypeTextPlain, []byte(`foo`), JSON(nil), PlainText(&text))
	if err != nil {
		t.Fatal(err)
	}
	if text != "foo" {
		t.Fatal("invalid value of text")
	}
}

func TestGetDefaultJSON400(t *testing.T) {
	err := get(http.StatusBadRequest, ContentTypeJSON, []byte(`{"foo": 42}`))
	if err == nil {
		t.Fatal("err is nil")
	}
	respErr := getResponseError(err)
	if respErr.Body != nil {
		if o, ok := respErr.Body.(map[string]interface{}); ok {
			if o["foo"] != nil {
				return
			}
		}
		t.Fatalf("invalid body: %#v", respErr.Body)
	}
	t.Fatal(respErr)
}

func TestGetDefaultText400(t *testing.T) {
	err := get(http.StatusBadRequest, ContentTypeTextPlain, []byte(`foo`))
	if err == nil {
		t.Fatal("err is nil")
	}
	respErr := getResponseError(err)
	if respErr.Body != nil {
		if respErr.Body == "foo" {
			return
		}
		t.Fatalf("invalid body: %#v", respErr.Body)
	}
	t.Fatal(respErr)
}
