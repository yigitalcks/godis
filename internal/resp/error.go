package resp

import "fmt"

type ErrTyp int

const (
	WrongFormat ErrTyp = iota
	EncodingError
	EmptyArray
	WrongCommand
	RequestTooLarge
	ConnectionError
)

var ErrPrefixMap = map[ErrTyp]string{
	WrongFormat:     "ERR",
	EncodingError:   "ERR",
	EmptyArray:      "ERR",
	WrongCommand:    "ERR",
	RequestTooLarge: "ERR",
	ConnectionError: "ERR",
}

var DefaultErrMsgMap = map[ErrTyp]string{
	WrongFormat:     "request is in invalid format",
	EncodingError:   "an error has occurred during encoding",
	EmptyArray:      "request is an empty array",
	WrongCommand:    "unknown command",
	RequestTooLarge: "max request size exceeded",
	ConnectionError: "connection error",
}

func (r *RespSimpleError) Error() string {
	return fmt.Sprintf("%s, %s", ErrPrefixMap[r.Prefix], r.Message)
}
