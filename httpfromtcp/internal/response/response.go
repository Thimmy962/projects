package response

import (
	"fmt"
	"http/internal/headers"
	"io"
)


type StatusCode int

const (
	success StatusCode = 200
	bad = 400
	server_error = 500
)


var responseEnum = map[StatusCode]string{
	200: "HTTP/1.1 200 OK\r\n",
	400: "HTTP/1.1 400 Bad Request\r\n",
	500: "HTTP/1.1 500 Internal Server Error\r\n",
}

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	_, err := w.Write([]byte(responseEnum[statusCode]))
	if err != nil {
		return err
	}
	return nil
}


func GetDefaultHeaders(contentLen int) headers.Headers {
	headers := make(headers.Headers)
	headers.ParseExistingFieldName("Content-Length", fmt.Sprintf("%d", contentLen))
	headers.ParseExistingFieldName("Connection", "Closed")
	headers.ParseExistingFieldName("Content-Type", "text/plain")
	return headers
}


func WriteHeaders(w io.Writer, headers headers.Headers) error {
	headerBytes := headers.Bytes()
	if _, err := w.Write(headerBytes); err != nil {
		return err
	}
	return nil
}
