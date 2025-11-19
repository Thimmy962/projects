package headers

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"
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
	h.ParseExistingFieldName(fieldName, lineValue)

	// number of bytes parsed index + 2 for the newline char
	return index + len(SEPARATOR), false, nil
}

func (h Headers) CheckHeader(key string) bool {
	key = strings.ToLower(key)
	_, ok := h[key]
	return ok
}

func( h Headers) Get(key string) string {
	key = strings.ToLower(key)
	return h[key]
}

// used if a particular header can have multiple values
func (h Headers) ParseExistingFieldName(fieldName, lineValue string) {
	// if the fieldName already exists
	if _, exist := h[fieldName]; exist {
		h[fieldName] = fmt.Sprintf("%s, %s", h[fieldName], lineValue)
	} else { 
		h[fieldName] = lineValue
	}
}

// converts the content of headers into bytes of field mapped to its value
func (h Headers) Bytes() []byte {
    var buf bytes.Buffer
    for k, v := range h {
        buf.WriteString(fmt.Sprintf("%s: %s", k, strings.TrimSuffix(v, "\r\n")))
        buf.WriteString("\r\n")
    }
    buf.WriteString("\r\n") // end of header section
    return buf.Bytes()
}

// checks the characters in the field name
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

// if the key exist assign a value to that key else create that key in the map and assign value to it.
func (h Headers) Set(key,value string) {
	h[key] = value
}

func (h Headers) Delete(key string) {
	delete(h, key)
}