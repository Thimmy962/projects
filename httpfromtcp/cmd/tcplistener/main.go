package main

import (
	"http/internal/request"
	"fmt"
	
	"net"
)

/*func getLinesChannel(conn io.ReadCloser) <-chan string {
	str_chan := make(chan string)
	go func () {
		defer conn.Close()
		defer close(str_chan)
		file := bufio.NewScanner(conn)
		var line string
		for file.Scan() {
			line = file.Text()
			str_chan <-line
		}
	}()
	return str_chan
}
*/

func main() {
	ln, err := net.Listen("tcp", ":42069"); if err != nil {
		panic(err)
	}
	defer ln.Close()
for {
    conn, err := ln.Accept()
    if err != nil {
        fmt.Println(err)
        continue
    }

    go func(conn net.Conn) {
        defer conn.Close()

        req, err := request.RequestFromReader(conn)
        if err != nil {
            fmt.Println("Error parsing request:", err)
            return
        }

        fmt.Printf("Request line:\n- Method: %s\n- Target: %s\n- Version: %s\n", 
            req.RequestLine.Method, req.RequestLine.RequestTarget, req.RequestLine.HttpVersion)
        fmt.Println("Headers:")
        for header := range req.Headers {
            fmt.Printf("- %s: %s\n", header, req.Headers[header])
        }

        if len(req.Body) > 0 {
            fmt.Printf("Body:\n%s\n", string(req.Body))
        }

        // ✅ respond so curl doesn’t block
        response := "HTTP/1.1 200 OK\r\n" +
            "Content-Type: text/plain\r\n" +
            "Content-Length: 12\r\n" +
            "Connection: close\r\n" +
            "\r\n" +
            "Hello World!\n"

        conn.Write([]byte(response))
    }(conn)
	}
}
