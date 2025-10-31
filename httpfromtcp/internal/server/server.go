package server

import (
	"fmt"
	"http/internal/request"
	"http/internal/response"
	"log"
	"net"
)

type Server struct {
	closed bool
	ln net.Listener
}

func Serve(port string) (*Server, error) {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalln(err.Error())
		return nil, err
	}
	
	serve := Server{false, ln}
	go serve.listen()
	return &serve, nil
}


func (s *Server) Close() error {
		s.closed = true
		return s.ln.Close()
}


func (s *Server) handle(connection net.Conn) {
	_, err := request.RequestFromReader(connection)
	if err != nil {
		fmt.Println(err.Error())
		s.Close()
		return
	}
	data := "Hello World!\n"
	err = response.WriteStatusLine(connection, 200)
	if err != nil {
		s.Close()
	}

	headers := response.GetDefaultHeaders(0)

	if err = response.WriteHeaders(connection, headers); err != nil {
		s.Close()
		return
	}
	
	connection.Write([]byte(data))
}


func (s *Server) listen() {
	for {
		if s.closed {
			break
		}
		conn, err := s.ln.Accept()
		if err != nil {
			fmt.Println(err.Error())
		}

		go func(conn net.Conn) {
			defer conn.Close()
			s.handle(conn)
		}(conn)
	}
}

