package main

import (
	"http/internal/headers"
	"http/internal/request"
	"http/internal/response"
	"http/internal/server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const port = "42069"



func Handle(w *response.Writer, req *request.Request) *server.HandlerError {
	headers := make(headers.Headers)
	handlerError := &server.HandlerError{}
	
	headers.ParseExistingFieldName("Connection", "Closed")
	
	handle := server.GetFunc(req.RequestLine.ParsedUrl)	

	handle(handlerError, &headers)
	



	err := w.WriteStatusLine(response.StatusCode(handlerError.Code()))
	if err != nil {
		server.ErrorWriting(err, handlerError)
		return handlerError
	}


	if err = w.WriteHeaders(headers); err != nil {
		server.ErrorWriting(err, handlerError)
		return handlerError
	}


	if _, err = w.WriteBody([]byte(handlerError.ErrMsg())); err != nil {
		server.ErrorWriting(err, handlerError)
		return handlerError
	}
	
	return handlerError
}
func main() {
	server, err := server.Serve(port, Handle)
	// server, err := server.Serve(port, handler.greet())
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}