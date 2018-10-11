package httpsimp

import (
	"fmt"
	"mime"
	"net/http"
)

/*
Parser matches and handles an http.Response.

To create a parser, use of the built-in parser functions like JSON,
PlainText, etc, or build a custom one using MakeParser.
*/
type Parser struct {
	ctype      string
	statusSpec StatusSpec
	retErr     bool
	parseBody  func(resp *http.Response) (interface{}, error)
}

/*
ParseOption is passed into MakeParser and built-in parser functions
to adjust which responses the parser matches and whether it
matches an error response.

You cannot define custom parser options.
*/
type ParseOption interface {
	applyToParser(m *Parser)
}

/*
MakeParser builds a parser wrapping the given parse function.

The parser starts out matching responses with the given content type
(which can be empty to match any response).

The provided options change the behavior of the parser and may
override the content type that it matches.
*/
func MakeParser(defaultCtype string, mopt []ParseOption, bodyParser func(resp *http.Response) (interface{}, error)) Parser {
	p := Parser{defaultCtype, Status2xx, false, bodyParser}
	for _, o := range mopt {
		o.applyToParser(&p)
	}
	return p
}

type matchOptionFunc func(m *Parser)

func (o matchOptionFunc) applyToParser(m *Parser) {
	o(m)
}

/*
ContentType causes the parser to only match responses with the given content type.
If an empty string is passed in, the parser will match any content type.
*/
func ContentType(ctype string) ParseOption {
	return matchOptionFunc(func(m *Parser) {
		m.ctype = ctype
	})
}

/*
ReturnError causes Do or Parse to return a non-nil error if this
parser matches. (The body is still parsed and handled.)
*/
func ReturnError() ParseOption {
	return returnError
}

var returnError ParseOption = matchOptionFunc(func(m *Parser) {
	m.retErr = true
})

func (s StatusSpec) applyToParser(m *Parser) {
	m.statusSpec = s
}

func parse(resp *http.Response, p Parser) (bool, error) {
	mediaType := resp.Header.Get("Content-Type")
	ctype, _, err := mime.ParseMediaType(mediaType)
	if err != nil {
		return false, fmt.Errorf("cannot parse Content-Type string %v", mediaType)
	}

	ctypeOK := (p.ctype == "" || ctype == p.ctype)
	statusOK := p.statusSpec.Matches(resp.StatusCode)
	if !ctypeOK || !statusOK {
		return false, &responseError{
			StatusCode:        resp.StatusCode,
			ContentType:       ctype,
			WantedContentType: p.ctype,
			ContentTypeOK:     ctypeOK,
			Body:              nil,
			DecodingError:     nil,
		}
	}

	body, bodyErr := p.parseBody(resp)
	if p.retErr || bodyErr != nil {
		return true, &responseError{
			StatusCode:        resp.StatusCode,
			ContentType:       ctype,
			WantedContentType: p.ctype,
			ContentTypeOK:     true,
			Body:              body,
			DecodingError:     bodyErr,
		}
	} else {
		return true, nil
	}
}

var fallbackParsers = []Parser{
	JSON(nil, Status4xx5xx, ReturnError()),
	PlainText(nil, Status4xx5xx, ContentType(ContentTypeTextPlain), ReturnError()),
	None(StatusAny, ReturnError()),
}

/*
Parse handles the HTTP response using of the provided parsers.
The first matching parser wins.

If no parsers match, some predefined fallback parsers are tried;
all of them cause a non-nil error to be returned.
*/
func Parse(resp *http.Response, parsers ...Parser) error {
	var firstErr error

	for _, p := range parsers {
		matched, err := parse(resp, p)
		if matched {
			return err
		}
		if firstErr == nil {
			firstErr = err
		}
	}

	for i, p := range fallbackParsers {
		matched, err := parse(resp, p)
		if matched {
			if i == len(fallbackParsers)-1 && err != nil {
				err = firstErr
			}
			return err
		}
	}

	return nil
}
