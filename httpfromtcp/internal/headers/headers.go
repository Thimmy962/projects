package headers

import (
	"fmt"
	"strings"
	"unicode"
	// "bytes"
)

type Headers map[string]string
var SEPARATOR = "\r\n"
var INVALID_HEADER_FIELDNAME = fmt.Errorf("invalid fieldname for http header")
var INVALID_FIELD_NAME = fmt.Errorf("invalid character in fieldname")

// parse the each headerline into a map[string]string map
func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	s := string(data)
	
	// Then look for CRLF elsewhere (a regular header line)
	index := strings.Index(s, SEPARATOR)
	switch index {
		case -1:
			return 0, false, nil
		case 0:
			return len(SEPARATOR), true, nil
	}


	data = data[:index] //create a new slice from the beginning to the char before the newline character

	sData := strings.TrimSpace(string(data)) //converts bytes into a string then trim all white spaces
	
	if strings.Contains(sData, " :") {
		return 0, false, INVALID_HEADER_FIELDNAME
	}


	spData := strings.SplitN(sData, ":", 2) // split on : into 2
	if len(spData) != 2 {
		return 0, false, fmt.Errorf("malformed header line")
	}


	fieldName := strings.ToLower(spData[0])
	// check if the strings in fieldName have invalid char
	if err := checkFieldName(fieldName); err != nil {
		return 0, false, err
	}

	lineValue := strings.TrimSpace(spData[1])
	h.parseExistingFieldName(fieldName, lineValue)

	// number of bytes parsed index + 2 for the newline char
	return index + len(SEPARATOR), false, nil
}

func( h Headers) Get(key string) (string, error) {
	key = strings.ToLower(key)
	if _, exist := h[key]; exist {
		return h[key], nil
	}
	return "Nothing", fmt.Errorf("this header does not exist")
}

func (h Headers) parseExistingFieldName(fieldName, lineValue string) {
	// if the fieldName already exists
	if _, exist := h[fieldName]; exist {
		h[fieldName] = fmt.Sprintf("%s, %s", h[fieldName], lineValue)
	} else { 
		h[fieldName] = lineValue
	}
}


func checkFieldName(fieldname string) error {
	for _, char := range fieldname {
		if unicode.IsDigit(char) || unicode.IsLetter(char) || specificSymbol(char) {
			continue
		}else {
			return INVALID_FIELD_NAME
		}
	}
	return nil
}


var allowedSymbols = map[rune]bool{
	'!': true, '#': true, '$': true, '%': true, '&': true, '\'': true,
	'*': true, '+': true, '-': true, '.': true, '^': true, '_': true,
	'`': true, '|': true, '~': true,
}

func specificSymbol(r rune) bool {
	return allowedSymbols[r]
}
