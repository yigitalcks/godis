package resp

import "bufio"

var (
	PrefixArray        = []byte{'*'}
	PrefixSimpleString = []byte{'+'}
	PrefixSimpleError  = []byte{'-'}
	PrefixInteger      = []byte{':'}
	PrefixBulkString   = []byte{'$'}

	CRLF = []byte{'\r', '\n'}
)

type RespType interface {
	Encode(w *bufio.Writer) error
}

type RespSimpleString struct {
	Value string
}

type RespSimpleError struct {
	Value string
}

type RespInteger struct {
	Value int
}

type RespBulkString struct {
	Value []byte
}

type RespArray struct {
	Value []RespType
}
