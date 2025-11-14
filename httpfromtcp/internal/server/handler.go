package server

import (
	"http/internal/request"
	"http/internal/response"
	"io"
)

type HandlerError struct {
	StatusCode int
	// errMsg is not neccessary an error just message to be passed
	ErrorMsg string
}



type Handler func(w io.Writer, req *request.Request) *HandlerError

// writes an handlerError to io.Writer
func (hErr *HandlerError)HandlerToWriter(w io.Writer) {
	response.WriteStatusLine(w, response.StatusCode(hErr.StatusCode))
	headers := response.GetDefaultHeaders(len(hErr.ErrorMsg))
	response.WriteHeaders(w, headers)
	w.Write([]byte(hErr.ErrorMsg))
}