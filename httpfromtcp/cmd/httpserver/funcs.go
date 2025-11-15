package main

import (
	"http/internal/response"
	"http/internal/server"
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
			headers.ParseExistingFieldName("Content-Type", "text/html")
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
			headers.ParseExistingFieldName("Content-Type", "text/html")
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
			headers.ParseExistingFieldName("Content-Type", "text/html")
			w.Headers = headers
			w.StatusCode = 200
			w.Body.Write([]byte(data))
			return nil	
}