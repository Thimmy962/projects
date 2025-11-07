package server

import (
	"io"
	"http/internal/request"
	"fmt"
)

type HandlerError struct {
	code int
	errorMsg string
}

// the constr
func HandlerErrorConstructor(code int, msg string) *HandlerError{
	return &HandlerError{code, msg}
}



type Handler func(w io.Writer, req *request.Request) *HandlerError


func Greet(w io.Writer, req *request.Request) *HandlerError {
	_,err := w.Write([]byte("Hello\n"))
	
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Error %d: %s", 500, err.Error())))
		return &HandlerError{500, err.Error()}
	}
	return nil
}