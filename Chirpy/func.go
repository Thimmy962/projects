package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"fmt"
)


// called on bad request
func ProcessingError(w http.ResponseWriter, code int, Err error) {
	w.Header().Set("Content-Type", "Apllication/json")
	w.WriteHeader(code)
	w.Write([]byte(Err.Error()))
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
	w.WriteHeader(code)
	buf.Write(jsonData)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", buf.Len()))
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(code)
	w.Write(buf.Bytes())

}


func respondWithJSONWithoutBuffer(w http.ResponseWriter, code int, payload interface{}) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(500)
		return
	}
	// convert jsonData to string from byte, pass it to a function that returns a string
	// convert the returned string to byte and svave in jsonData 
	jsonData = []byte(profaneFUnc(string(jsonData)))
	w.WriteHeader(code)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(code)

}



func ValidateChirp(w http.ResponseWriter, req *http.Request) {
	var buf bytes.Buffer
	n, err := buf.ReadFrom(req.Body)
	// if there is an error in the writing
	if err != nil {
		ProcessingError(w, 500, fmt.Errorf("something went wrong"))
		return
	}
	// if the number of bytes read(n) is > 140 
	if n > 140 {
		ProcessingError(w, 400, fmt.Errorf("chirp too long"))
		return
	}
	chirp := Chirp{Body: buf.String()}
	respondWithJSON(w, 200, buf, chirp)
}