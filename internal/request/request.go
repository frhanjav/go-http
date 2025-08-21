package request

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
}

var Err_bad = fmt.Errorf("bad error")
var SEPARATOR = "\r\n"

func parseRequestLines(b string) (*RequestLine, string, error) {
	idx := strings.Index(b, SEPARATOR)
	
	if idx == -1 {
		return nil, b, nil
	}

	startLine := b[:idx]
	restOfMessage := b[idx+len(SEPARATOR):] 
	
	parts := strings.Split(startLine, " ")

	if len(parts) != 3 {
		return nil, restOfMessage, Err_bad
	}

	httpParts := strings.Split(parts[2], "/")
	if len(parts) != 2 || httpParts[0] != "HTTP" || httpParts[1] != "1.1" {
		return nil, restOfMessage, Err_bad
	}

	rl := &RequestLine{
	HttpVersion:httpParts[1],
	RequestTarget:parts[1],
	Method:parts[0],
	}

	return rl, restOfMessage, nil 
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	data ,err := io.ReadAll(reader)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("unable to readall"), err,)
	}

	str := string(data)

	rl, _, err := parseRequestLines(str)
	if err != nil {
		return nil, err
	}

	return &Request{
		RequestLine: *rl,
	}, err
}