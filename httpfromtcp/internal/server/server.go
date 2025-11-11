package server

import (
	"bytes"
	"fmt"
	"http/internal/request"
	"http/internal/response"
	"log"
	"net"
)

type Server struct {
	closed bool
	ln net.Listener
	handlerFunction Handler
}

func Serve(port string, handler Handler) (*Server, error) {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalln(err.Error())
		return nil, err
	}
	
	serve := &Server{false, ln, handler}
	go serve.listen()
	return serve, nil
}


func (s *Server) Close() error {
		s.closed = true
		return s.ln.Close()
}


func (s *Server) handle(connection net.Conn) {
	defer connection.Close()
	parsedRequest, err := request.RequestFromReader(connection)
	if err != nil {
		s.Close()
		return
	}
	
	var buf bytes.Buffer
	writer := response.InitWriter(&buf)

	s.handlerFunction(&writer, parsedRequest)



	if _, err = connection.Write(buf.Bytes()); err != nil {
			s.Close()
			return
		}
}


func (s *Server) listen() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			if s.closed {
				return
			}
			fmt.Println(err.Error())
		}

		go s.handle(conn)
	}
}

