package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
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


// handles both writing a single struct or a list of structs to responseWriter
// all thanks to NewEncoder
func respondWithJSON(w http.ResponseWriter, code int, payload any) {
	var buf bytes.Buffer

	err := json.NewEncoder(&buf).Encode(payload)
	if err != nil {
		log.Println(err.Error())
		ProcessingError(w, 500, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", buf.Len()))
	w.WriteHeader(200)

	err = json.NewEncoder(w).Encode(payload)
	if err != nil {

	}
}


// deprecated
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
	respondWithJSON(w, 200, chirp)
}
