package main

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"http/internal/request"
	"http/internal/response"
	"http/internal/server"
	"io"
	"net/http"
	"os"
	// "path/filepath"
)


func yourproblem(w *response.Writer) *server.HandlerError {
	data:= "<html>"+
  					"<head>"+
    					"<title>400 Bad Request</title>"+
  					"</head>"+
  					"<body>"+
    					"<h1>Bad Request</h1>"+
    					"<p>Your request honestly kinda sucked.</p>"+
  					"</body>"+
				"</html>"
			headers := response.GetDefaultHeaders(len(data))
			headers.Set("Content-Type", "text/html")
			w.Headers = headers
			w.StatusCode = 400
			w.Body.Write([]byte(data))
		return &server.HandlerError{StatusCode: int(w.StatusCode), ErrorMsg: data}
}

func myproblem(w *response.Writer) *server.HandlerError {
		data:= "<html>"+
  					"<head>"+
    					"<title>500 Internal Server Error</title>"+
  					"</head>"+
  					"<body>"+
    					"<h1>Internal Server Error</h1>"+
    					"<p>Okay, you know what? This one is on me</p>"+
  					"</body>"+
				"</html>"
			headers := response.GetDefaultHeaders(len(data))
			headers.Set("Content-Type", "text/html")
			w.Headers = headers
			w.StatusCode = 500
		return &server.HandlerError{StatusCode: int(w.StatusCode), ErrorMsg: data}
}


func Default(w *response.Writer) *server.HandlerError {
		data:= "<html>"+
  					"<head>"+
    					"<title>200 Ok</title>"+
  					"</head>"+
  					"<body>"+
    					"<h1>Success!</h1>"+
    					"<p>Your request was an absolute banger</p>"+
  					"</body>"+
				"</html>"
			headers := response.GetDefaultHeaders(len(data))
			headers.Set("Content-Type", "text/html")
			w.Headers = headers
			w.StatusCode = 200
			w.Body.Write([]byte(data))
			return nil	
}


func stream(w *response.Writer, req *request.Request) *server.HandlerError {
	url := req.RequestLine.ParsedUrl
	w.Headers = response.GetStreamDefaultHeaders()
	res, err := http.Get("http://localhost:8000/" + url[1] + "/" + url[2])
	if err != nil {
		return myproblem(w)
	}
	w.StatusCode = 200
	w.WriteStatusLine(200)
	w.WriteHeaders()
	bufLen := 1024
	buf := make([]byte, bufLen)

	dataLen := 0
	hash256 := sha256.New()

	for {
		n, rErr := res.Body.Read(buf); if rErr != nil {
			if !errors.Is(rErr, io.EOF) {
				yourproblem(w)
			}
		}

		if n > 0 {
			hash256.Write(buf)
			length, _ := w.WriteChunkBody(buf)
			dataLen += length
		} else {
			w.WriteChunkedBodyDone()
			break
		}
	}
	finalhasH := hash256.Sum(nil)
	clear(w.Headers)

	w.Headers.ParseExistingFieldName("X-Content-SHA256", fmt.Sprintf("%x", finalhasH))
	w.Headers.ParseExistingFieldName("X-Content-Length", fmt.Sprintf("%d", dataLen))

	w.WriteTrailers(w.Headers)
	
	return nil
}



func video(w *response.Writer, _ *request.Request) *server.HandlerError {
	w.Headers = response.GetStreamDefaultHeaders()
	w.Headers.Set("Content-Type", "video/mp4")
	// w.Headers.Set("Content-Disposition", "attachment; filename=abc.mp4")
	w.Headers.ParseExistingFieldName("Content-Type", "video/webm")
	w.Headers.ParseExistingFieldName("Content-Type", "video/ogg")
	w.Headers.Set("Accept-Ranges","bytes")

	w.StatusCode = 200

	bufLen := 1024
	buf := make([]byte, bufLen)

	w.WriteStatusLine(200)
	w.WriteHeaders()

	filepath := "/home/timileyin/projects/httpfromtcp/assets/vim.mp4"

	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println(err.Error())
		return myproblem(w)
	}

	defer file.Close()

	

	dataLen := 0
	hash256 := sha256.New()

	for {
		n, rErr := file.Read(buf); if rErr != nil {
			if !errors.Is(rErr, io.EOF) {
				yourproblem(w)
			}
		}

		if n > 0 {
			hash256.Write(buf)
			length, _ := w.WriteChunkBody(buf)
			dataLen += length
		} else {
			w.WriteChunkedBodyDone()
			break
		}
	}
	finalhasH := hash256.Sum(nil)
	clear(w.Headers)

	w.Headers.Set("X-Content-SHA256", fmt.Sprintf("%x", finalhasH))
	w.Headers.Set("X-Content-Length", fmt.Sprintf("%d", dataLen))

	w.WriteTrailers(w.Headers)
	
	return nil
}