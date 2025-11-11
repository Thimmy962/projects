package response

import (
	"fmt"
	"http/internal/headers"
	"io"
	"bytes"
)


type StatusCode int

const (
	success StatusCode = 200
	bad = 400
	server_error = 500
)

type Writer struct {
	write io.Writer
	status int
}

func InitWriter(writer *bytes.Buffer) Writer{
	return Writer{writer, 0}
}

var responseEnum = map[StatusCode]string{
	200: "HTTP/1.1 200 OK\r\n",
	400: "HTTP/1.1 400 Bad Request\r\n",
	500: "HTTP/1.1 500 Internal Server Error\r\n",
}


func GetDefaultHeaders(contentLen int) headers.Headers {
	headers := make(headers.Headers)
	headers.ParseExistingFieldName("Content-Length", fmt.Sprintf("%d", contentLen))
	headers.ParseExistingFieldName("Connection", "Closed")
	headers.ParseExistingFieldName("Content-Type", "text/plain")
	return headers
}


func (w *Writer)WriteStatusLine(statusCode StatusCode) error {
	_, err := w.write.Write([]byte(responseEnum[statusCode]))
	w.status = 1
	return err
}

func (w *Writer)WriteHeaders(headers headers.Headers) error {
	if w.status != 1 {
		panic("YOU HAVE TO CALL THE WriteStatusLine method of the Writer type first")
	}
	// Bytes method was implemented in the headers package to turn headers into bytes
	headerBytes := headers.Bytes()
	_, err := w.write.Write(headerBytes)
	w.status = 2
	return err
}


func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.status != 2 {
		panic("YOU HAVE TO CALL THE WriteHeaders method of the Writer type first")
	}
	return w.write.Write(p)
}


func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	return w.write.Write(p)
}

// func (w *Writer) WriteChunkedBodyDone() (int, error)