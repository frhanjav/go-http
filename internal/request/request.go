package request

import (
	"bytes"
	"fmt"
	"io"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type ParserState string

const (
	StateInit ParserState = "init"
	StateDone ParserState = "done"
	StateError ParserState = "error"
)


type Request struct {
	RequestLine RequestLine
	state ParserState
}

var Err_bad = fmt.Errorf("bad error")
var Err_unsupported_http_ver = fmt.Errorf(("unsupported http version"))
var Err_request_in_error_state = fmt.Errorf("request in error state")
var SEPARATOR = []byte("\r\n")

func newRequest() *Request {
	return &Request {
		state: StateInit,
	}
}

func (r *Request) parse(data []byte) (int, error) {

	read := 0
	outer:
	for {
		switch r.state {
		case StateError:
			return 0, Err_request_in_error_state
		case StateInit:
			rl, n, err := parseRequestLines(data[read:])
			if err != nil {
				return 0, err
			}
			if n == 0 {
				break outer
			}
			
			r.RequestLine = *rl
			read += n

			r.state = StateDone
			
		case StateDone:
			break outer
		}
	}
	return read, nil
}

func (r *Request) done() bool {
	return r.state == StateDone || r.state == StateError
}

// here int is the no. of bytes that have been read
func parseRequestLines(b []byte) (*RequestLine, int, error) {
	idx := bytes.Index(b, SEPARATOR)
	
	if idx == -1 {
		return nil, 0, nil
	}

	startLine := b[:idx]
	read := idx+len(SEPARATOR)
	
	parts := bytes.Split(startLine, []byte(" "))

	if len(parts) != 3 {
		return nil, 0, Err_bad
	}

	httpParts := bytes.Split(parts[2], []byte("/"))
	if len(parts) != 2 || string(httpParts[0]) != "HTTP" || string(httpParts[1]) != "1.1" {
		return nil, 0, Err_unsupported_http_ver
	}

	rl := &RequestLine{
	HttpVersion:   string(httpParts[1]),
	RequestTarget: string(parts[1]),
	Method:        string(parts[0]),
	}

	return rl, read, nil 
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := newRequest()
	buf := make([]byte, 1024)
	bufIdx := 0
	for !request.done() {
		n, err := reader.Read(buf[bufIdx:])
		if err != nil {
			return nil, err
		}

		bufIdx += n
		readN, err := request.parse(buf[:bufIdx])
		if err != nil {
			return nil, err
		}

		copy (buf, buf[readN:bufIdx])
		bufIdx -= readN
	}

	return request, nil
}