package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"fmt"
)

// called on internal server error for application Json
func Err_500ApplicationJson(w http.ResponseWriter, Err string) {
	w.Header().Set("Content-Type", "Aplication/json")


	resBody := Error{Err: Err}

	jsonData, err := json.Marshal(resBody)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(500)
	w.Write([]byte(jsonData))
}

// called on bad request
func Err_400ApplicationJson(w http.ResponseWriter, Err string) {
	w.Header().Set("Content-Type", "Apllication/json")
	resBody:= Error{Err: Err}

	jsonData, err :=  json.Marshal(resBody)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(400)
	w.Write([]byte(jsonData))
}


func profaneFUnc(str string) string {
	for _, word := range profane {
		str = cutAndJoin(word, str)
	}
	return str
}


func cutAndJoin(subStr, str string) string {
	// how non overlapping times is substr present in str
	count := strings.Count(str, subStr)
	// replace substr count times with **** in str
	return strings.Replace(str, subStr, "****", count)
}



func respondWithJSON(w http.ResponseWriter, code int, buf bytes.Buffer, payload interface{}) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(500)
		return
	}
	buf.Reset()
	// convert jsonData to string from byte, pass it to a function that returns a string
	// convert the returned string to byte and svave in jsonData 
	jsonData = []byte(profaneFUnc(string(jsonData)))
	buf.Write(jsonData)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", buf.Len()))
	
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(code)
	w.Write(buf.Bytes())

}