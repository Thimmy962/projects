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



type Handler func(w *response.Writer, req *request.Request) *HandlerError

// writes an handlerError to io.Writer
func HandlerToWriter(w io.Writer, p []byte) {
	w.Write(p)
}