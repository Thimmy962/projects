package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"http/internal/server"
	"http/internal/request"
	"io"
)

const port = "42069"

// var handle server.Handler = server.Greet

func Handle(w io.Writer, req *request.Request) *server.HandlerError {
	switch req.RequestLine.RequestTarget {
		case "yourproblem":
			w.Write([]byte("Your problem is not my problem\n"))
			return server.HandlerErrorConstructor(400, "Your problem is not my problem\n")
		case "myproblem":
			w.Write([]byte("Woopsie, my bad\n"))
			return server.HandlerErrorConstructor(500, "Woopsie, my bad\n")
		default:
			w.Write([]byte("All good frfr\n"))
			return server.HandlerErrorConstructor(200,  "All good, frfr\n")
	}
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