package main

import (
	"http/internal/request"
	"http/internal/response"
	"http/internal/server"
	"log"
	"os"
	"os/signal"
	"regexp"
	"syscall"
)

const port = "42069"


func handle(w *response.Writer, req *request.Request) *server.HandlerError{
	
	switch{
	case regexp.MustCompile(`^.yourproblem`).MatchString(req.RequestLine.RequestTarget):
		return yourproblem(w)
	case  regexp.MustCompile(`^.myproblem`).MatchString(req.RequestLine.RequestTarget):
		return myproblem(w)
	default:
		return Default(w)
	}
}


func main() {
	server, err := server.Serve(port, handle)
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