package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
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