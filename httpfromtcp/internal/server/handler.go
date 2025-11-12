package server

import (
	"fmt"
	"http/internal/request"
	"http/internal/response"
	"http/internal/headers"
	"io"
)

type HandlerError struct {
	statusCode int
	errorMsg string
}

// modifies the value of HandlerError
func (handlerError *HandlerError) HandlerErrorModifier(code int, msg string) {
	handlerError.statusCode = code
	handlerError.errorMsg = msg
}


func (handlerError *HandlerError)ErrMsg() string {
	return handlerError.errorMsg
}

func (handlerError *HandlerError)Code() int {
	return handlerError.statusCode
}



type Handler func(w *response.Writer, req *request.Request) *HandlerError


func Greet(w io.Writer, req *request.Request) *HandlerError {
	_,err := w.Write([]byte("Hello\n"))
	
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Error %d: %s", 500, err.Error())))
		return &HandlerError{500, err.Error()}
	}
	return nil
}



// func Proxyhandler(w *response.Writer, req *request.Request) *HandlerError {
	// path := req.RequestLine.RequestTargetreq *request.Request) *HandlerError {
// }



// called when writing statusline, or header or body to buffer failed
func ErrorWriting(err error, handlerError *HandlerError) {
	data := "<html>" +
				"<head>"+
    				"<title>500 Internal Server Error</title>"+
  				"</head>"+
  				"<body>"+
    				"<h1>Internal Server Error</h1>"+
    				"<p>"+ err.Error() +"</p>"+
  				"</body>"+
			"</html>"
	handlerError.HandlerErrorModifier(500, data)
}

// if path is your problem
func yourproblem(err *HandlerError, header *headers.Headers) {
	data := "<html>" +
  				"<head>"+
    				"<title>400 Bad Request</title>"+
				"</head>"+
				"<body>"+
					"<h1>Bad Request</h1>"+
					"<p>Your request honestly kinda sucked.</p>"+
  				"</body>"+
			"</html>"
	err.HandlerErrorModifier(400, data)
	header.ParseExistingFieldName("Content-Type", "text/html")
	header.ParseExistingFieldName("content-Length", fmt.Sprintf("%d", len(data)))
}

// if path is my problem
func myproblem(err *HandlerError, header *headers.Headers) {
	data := "<html>" +
  				"<head>"+
    				"<title>500 Internal Server Error</title>"+
  				"</head>"+
  				"<body>"+
					"<h1>Internal Server Error</h1>"+
    				"<p>Okay, you know what? This one is on me.</p>"+
				"</body>"+
			"</html>"
	err.HandlerErrorModifier(500, data)
	header.ParseExistingFieldName("content-Length", fmt.Sprintf("%d", len(data)))
	header.ParseExistingFieldName("Content-Type", "text/html")
}

func Default(err *HandlerError, header *headers.Headers) {
	
	data := "<html>" +
				"<head>"+
					"<title>200 Ok</title>"+
  				"</head>"+
				"<body>"+
					"<h1 color: \"red\">Success</h1>"+
					"<p>Your request was an absolute banger.</p>"+
  				"</body>"+
			"</html>"
	err.HandlerErrorModifier(200, data)
	header.ParseExistingFieldName("Content-Type", "text/html")
	header.ParseExistingFieldName("content-Length", fmt.Sprintf("%d", len(data)))
}



func handleHttpBin(err *HandlerError, header *headers.Headers) {
	data := "<h1>hello World</h1>"
	err.HandlerErrorModifier(200, data)
	header.ParseExistingFieldName("Content-Type", "text/html")
	// header.ParseExistingFieldName("Transfer-Encoding", "chunked")
	header.ParseExistingFieldName("Transfer-Encoding", "chunked")
}


func GetFunc(paths []string) func(*HandlerError, *headers.Headers) {
	fmt.Println(paths)
	first := paths[0]

	switch first {
	case "yourproblem":
		return yourproblem
	
	case "myproblem":
		return myproblem
	case "httpbin":
		return handleHttpBin
	case "ok":
		return Default
	default:
		return yourproblem
	}
}