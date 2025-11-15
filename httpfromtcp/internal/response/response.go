package response

import (
	"fmt"
	"http/internal/headers"
	"net"
	"io"
)


type StatusCode int

const (
	Success StatusCode = 200
	Bad = 400
	Server_error = 500
)

type Writer struct {
	Connection net.Conn
	Body io.Writer
	Headers headers.Headers
	StatusCode StatusCode
}


func (write *Writer) Write(p []byte) (int, error) {
	return write.Connection.Write(p)
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
	headers.ParseExistingFieldName("Server", "NGINX")
	return headers
}

func WriteHeaders(w io.Writer, header headers.Headers) error{
	_, err := w.Write(header.Bytes())
	return err
}


func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	_, err := w.Write([]byte(responseEnum[statusCode]))
	return err
}


func (w *Writer) WriteStatusLine(statusCode StatusCode) error{
	return WriteStatusLine(w.Connection, statusCode)
}


func (w *Writer) WriteHeaders(headers headers.Headers) error {
	return WriteHeaders(w.Connection, headers)
}


func (w *Writer) WriteBody(p []byte) (int, error) {
	return w.Connection.Write(p)
}