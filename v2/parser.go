package httpsimp

import (
	"fmt"
	"mime"
	"net/http"
)

type Parser struct {
	ctype      string
	statusSpec StatusSpec
	retErr     bool
	parseBody  func(resp *http.Response) (interface{}, error)
}

type ParseOption interface {
	applyToParser(m *Parser)
}

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

func ContentType(ctype string) ParseOption {
	return matchOptionFunc(func(m *Parser) {
		m.ctype = ctype
	})
}

var ReturnError ParseOption = matchOptionFunc(func(m *Parser) {
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

var defaultParsers = []Parser{
	JSON(nil, Status4xx5xx, ReturnError),
	PlainText(nil, Status4xx5xx, ContentType(ContentTypeTextPlain), ReturnError),
	None(StatusAny, ReturnError),
}

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

	for i, p := range defaultParsers {
		matched, err := parse(resp, p)
		if matched {
			if i == len(defaultParsers)-1 && err != nil {
				err = firstErr
			}
			return err
		}
	}

	return nil
}
