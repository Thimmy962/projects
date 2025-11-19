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
	Bad StatusCode = 400
	Server_error StatusCode = 500
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

// These headers are set by default. 
// Headers struct has a Delete method that takes a key and remove it from the headers.Header map
// You can as well use the standard library delete(<Headers>, key) to delete
func GetDefaultHeaders(contentLen int) headers.Headers {
	headers := make(headers.Headers)
	headers.Set("Content-Length", fmt.Sprintf("%d", contentLen))
	headers.Set("Connection", "Closed")
	headers.Set("Content-Type", "text/plain")
	headers.Set("Server", "NGINX")
	return headers
}

// Same as GetDefaultHeaders just that 
func GetStreamDefaultHeaders() headers.Headers {
	headers := make(headers.Headers)
	headers.Set("Connection", "Closed")
	headers.Set("Transfer-Encoding", "chunked")
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


func (w *Writer) WriteHeaders() error {
	return WriteHeaders(w.Connection, w.Headers)
}


func (w *Writer) WriteBody(p []byte) (int, error) {
	return w.Connection.Write(p)
}


func (w *Writer) WriteChunkBody(p []byte) (int, error) {
	dataLen := len(p) // length of data in write to connection
	data := fmt.Sprintf("%x\r\n%s\r\n", dataLen, string(p)) // could have done all these in one line  but chose 3 for reading in the future
	w.Connection.Write([]byte(data))
	return len(p), nil
}


func (w *Writer) WriteChunkedBodyDone() (int, error) {
	return w.Connection.Write([]byte("0\r\n\r\n"))
}



func (w *Writer) WriteTrailers(h headers.Headers) error {
	w.WriteHeaders()
	return nil
}