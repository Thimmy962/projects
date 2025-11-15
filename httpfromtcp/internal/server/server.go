package server

import (
	"bytes"
	"fmt"
	"http/internal/headers"
	"http/internal/request"
	"http/internal/response"
	"log"
	"net"
	"net/url"
)

type Server struct {
	closed bool
	ln net.Listener
	handlerFunction Handler
}

func Serve(port string, handle Handler) (*Server, error) {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalln(err.Error())
		return nil, err
	}
	
	serve := &Server{false, ln, handle}
	go serve.listen()
	return serve, nil
}


func (s *Server) Close() error {
		s.closed = true
		return s.ln.Close()
}


func (s *Server) handle(connection net.Conn) {
	defer connection.Close()
	req, err := request.RequestFromReader(connection)
	if err != nil {
		response.WriteStatusLine(connection, 400)
		data := "Request or header badly formed"
		headers := headers.Headers{}
		headers.ParseExistingFieldName("Content-Type", "text/plain")
		headers.ParseExistingFieldName("Content-Length", fmt.Sprintf("%d",len(data)))
		response.WriteHeaders(connection, headers)
		connection.Write([]byte(data))
		return
	}

	url.Parse(req.RequestLine.RequestTarget)

	var buf bytes.Buffer
	writer := response.Writer{Connection: connection, Body: &buf}
	hErr := s.handlerFunction(&writer, req); 
	if hErr != nil {
		writer.WriteStatusLine(response.StatusCode(hErr.StatusCode))
		writer.WriteHeaders(writer.Headers)
		writer.WriteBody([]byte(hErr.ErrorMsg))
		return
	}


	writer.WriteStatusLine(response.StatusCode(200))
	writer.WriteHeaders(writer.Headers)
	
	_, err = writer.WriteBody(buf.Bytes())
	if err != nil {
		connection.Write([]byte(err.Error()))
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

