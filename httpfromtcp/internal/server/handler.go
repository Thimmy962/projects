package server

import (
	"fmt"
	"http/internal/response"
	"http/internal/request"
	"io"
)

type HandlerError struct {
	statusCode int
	errorMsg string
}

// modifies the value of HandlerError
func (handlerError *HandlerError) HandlerErrorModifier(code int, msg string) {
	handlerError.statusCode = code
	handlerError.errorMsg = msg
}


func (handlerError *HandlerError)ErrMsg() string {
	return handlerError.errorMsg
}

func (handlerError *HandlerError)Code() int {
	return handlerError.statusCode
}



type Handler func(w *response.Writer, req *request.Request) *HandlerError


func Greet(w io.Writer, req *request.Request) *HandlerError {
	_,err := w.Write([]byte("Hello\n"))
	
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Error %d: %s", 500, err.Error())))
		return &HandlerError{500, err.Error()}
	}
	return nil
}



// func Proxyhandler(w *response.Writer, req *request.Request) *HandlerError {
	// path := req.RequestLine.RequestTargetreq *request.Request) *HandlerError {
// }