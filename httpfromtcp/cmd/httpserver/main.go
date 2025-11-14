package main

import (
	// "http/internal/headers"
	"http/internal/request"
	"http/internal/server"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const port = "42069"


func handle(w io.Writer, req *request.Request) *server.HandlerError{
	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		return &server.HandlerError{400, "Your problem is not my problem\n"}
	case "/myproblem":
		return &server.HandlerError{500, "Woopsie, my bad\n"}
	default:
		data := "All good, frfr\n"
		w.Write([]byte(data))
		return nil
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