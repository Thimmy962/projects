package request

import (
	"bytes"
	"errors"
	"fmt"
	"http/internal/headers"
	"io"
	"strconv"
	"strings"
)

type internal int

const (
	initialized internal = iota
	parsingHeader
	parsingBody
	done
)

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	state       internal
	Body        []byte
}

var enum = map[internal]string{
	0: "initialized",
	1: "parsingHeader",
	2: "parsingBody",
	3: "done",
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

var MALFORMED_HTTP_VERSION = fmt.Errorf("malformed http version")
var MALFORMED_REQUEST_LINE = fmt.Errorf("malformed request line")

var SEPARATOR = "\r\n"

func RequestFromReader(r io.Reader) (*Request, error) {
	var buf bytes.Buffer
	bufLen := 10
	tmp := make([]byte, bufLen)
	request := Request{
		state:   initialized,
		Headers: make(headers.Headers),
	}

	for {

		n, err := r.Read(tmp)

		if n > 0 {
			buf.Write(tmp[:n])
		}
		
		// Always try to parse whatever is in buffer
		for buf.Len() > 0 {
			bytesParsed, perr := request.parse(buf.Bytes())
			if perr != nil {
				fmt.Println(perr.Error())
				return nil, perr
			}

			if bytesParsed > 0 {
				adjustingBuffer(bytesParsed, &buf)
			} else {
				break
			}

			if request.state == done {
				break
			}

		}

		// Check if parsing is fully done
		if request.state == done {
			break
		}


		if errors.Is(err, io.EOF){
			if buf.Len() == 0 {
				break
			}

			// There's still data left to parse even though EOF hit
			// → loop again to process that last chunk
			continue
		}
		if err != nil {
			return  nil, err
		}
	}

	return &request, nil
}

func (r *Request) parse(data []byte) (int, error) {
	switch r.state {
	case initialized:
		RLine, bytesRead, err := parseRequestLine(string(data))
		if RLine == nil { //if there is no error but nothing has been parsed into RequestLine
			return bytesRead, err
		}
		r.RequestLine = *RLine
		r.state = parsingHeader
		// fmt.Println("Done Parsing requestline")
		return bytesRead, err

	case parsingHeader:
		bytesRead, doneReading, err := r.Headers.Parse(data)
		if doneReading{
			r.state = parsingBody
		}
		return bytesRead, err

	case parsingBody:
		var content_length int
		value, err := r.Headers.Get("Content-Length")
		if err != nil {
			// err means no content-length in the header
			r.state = done
			return 0, nil
		}
		
		content_length64, intParsingErr := strconv.ParseInt(value, 10, 64)
		if intParsingErr != nil {
			return 0, intParsingErr
		}

		content_length = int(content_length64)
		if content_length < 1 {
			return 0, fmt.Errorf("content-length can not be less that 1")
		}


		dataLen := len(data)
		bodyLen := len(r.Body)

		// Body cannot exceed Content-Length
		if bodyLen > content_length {
			return 0, fmt.Errorf("excess body: already received %d bytes, expected %d", bodyLen, content_length)
		}

		// If new data was read, append it to the body
		if dataLen > 0 {
			r.Body = append(r.Body, data...)
			return dataLen, nil
		}
		
		if bodyLen < content_length {
				return 0, fmt.Errorf("incomplete body: received %d bytes, expected %d", bodyLen, content_length)
			}

			// perfect match - done parsing
		r.state = done

		return dataLen, nil
	}
	return -1, fmt.Errorf("error: trying to read data in a done state")
}

func parseRequestLine(line string) (*RequestLine, int, error) {
	/*
		reads line until the delimiter is met
		if there is no delimiter that means the entire line has not been read return 0
	*/
	index := strings.Index(line, SEPARATOR)
	if index == -1 {
		return nil, 0, nil
	}

	line = line[:index]

	// split the string option on space(" ")
	tokens := strings.Split(line, " ")
	if len(tokens) < 3 || !strings.Contains(tokens[2], "/") {
		return nil, 0, MALFORMED_REQUEST_LINE
	}

	httpVersionTokens := strings.Split(tokens[2], "/")
	if len(httpVersionTokens) != 2 || httpVersionTokens[0] != "HTTP" || httpVersionTokens[1] != "1.1" {
		return nil, 0, MALFORMED_HTTP_VERSION
	}

	return &RequestLine{
			HttpVersion: httpVersionTokens[1], RequestTarget: tokens[1], Method: tokens[0]},
		index + len(SEPARATOR), nil
}

func adjustingBuffer(bytesParsed int, buf *bytes.Buffer) {
	remaining := buf.Bytes()[bytesParsed:]
	buf.Reset()
	buf.Write(remaining)
}
