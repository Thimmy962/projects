package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func NewHeaders() Headers {
	return Headers{}
}

func TestRequestLineParse(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Invalid character in header
	headers = NewHeaders()
	data = []byte("H©st: localhost:42069\r\n\r\n")
	_, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.False(t, done)

	headers = NewHeaders()
	data = []byte("Host:localhost:8080\r\n\r\n")
	_, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.False(t, done)

	headers = NewHeaders()
	data = []byte("Host:localhost:8000\r\n\r\nServer: SimpleHTTP/0.6 /3.12.3\r\n\r\nHello")
	_, done, err = headers.Parse(data)

	require.NoError(t, err)
	assert.False(t, done)

}
