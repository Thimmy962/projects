package request

import (
	"bytes"
	"errors"
	"fmt"
	"http/internal/headers"
	"io"
	"strconv"
	"strings"
	"time"
	"net"
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



type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
	// This is not part of the question, but I added it any way so that I can use it when handling request and not parse again
	ParsedUrl []string
}


func (s *RequestLine) ParseUrl() []string {
	path := strings.TrimLeft(s.RequestTarget, "/")
	return strings.Split(path, "/")
}

var MALFORMED_HTTP_VERSION = fmt.Errorf("malformed http version")
var MALFORMED_REQUEST_LINE = fmt.Errorf("malformed request line")
var CONTENT_LENGTH_LENGTH_NOT_VALID = fmt.Errorf("the content length could not be parsed as it contains none digits")

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
		// If this is a network connection, and we're currently parsing the body,
		// set a short read deadline to allow graceful EOF detection.
		if request.state == parsingBody {
			if conn, ok := r.(net.Conn); ok {
				conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
			}
		}

		n, err := r.Read(tmp)

		if n > 0 {
			buf.Write(tmp[:n])
		} else {
			_, err := request.parse(buf.Bytes())
			if err != nil {
				return nil, err
			}
		}
		
		// Always try to parse whatever is in buffer
		for buf.Len() > 0 {
			bytesParsed, perr := request.parse(buf.Bytes())
			if perr != nil {
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
		// when the requesline has been gotten, parse the the target path
		r.RequestLine.ParsedUrl = r.RequestLine.ParseUrl()
		return bytesRead, err

	case parsingHeader:
		bytesRead, doneReading, err := r.Headers.Parse(data)
		if doneReading{
			if !r.Headers.CheckHeader("coNTent-length"){
				r.state = done
			} else {
				r.state = parsingBody
			}
		}
		return bytesRead, err

	case parsingBody:
		
		content_length, err := r.contentLengthParser()
		if err != nil {
			return 0, err
		}
		dataLen := len(data)
		bodyLen := len(r.Body)

		return r.concatData(dataLen, bodyLen, content_length, data)
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


func (r *Request) contentLengthParser() (int, error) {
	var content_length int
	value := r.Headers.Get("Content-Length")
	content_length64, intParsingErr := strconv.ParseInt(value, 10, 64)
	if intParsingErr != nil {
		return 0, fmt.Errorf("%s --- %s",intParsingErr,  CONTENT_LENGTH_LENGTH_NOT_VALID)
	}

	content_length = int(content_length64)
	if content_length < 1 {
		return 0, fmt.Errorf("content-length can not be less that 1")
	}
	return content_length, nil
}

/*
	Append a new string of data to the body
*/

func (r *Request)concatData(dataLen, bodyLen, contentLength int, data []byte) (int, error) {
		if bodyLen > contentLength {
			return 0, fmt.Errorf("excess body: already received %d bytes, expected %d", bodyLen, contentLength)
		}

		// If new data was read, append it to the body
		if dataLen > 0 {
			r.Body = append(r.Body, data...)
			return dataLen, nil
		}

		if bodyLen < contentLength {
				return 0, fmt.Errorf("incomplete body: received %d bytes, expected %d", bodyLen, contentLength)
		}

			// perfect match - done parsing
		r.state = done

		return dataLen, nil
}