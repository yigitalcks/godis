package resp

import "fmt"

type ErrTyp int

const (
	WrongFormat ErrTyp = iota
	EncodingError
	EmptyArray
	WrongCommand
)

var ErrPrefixMap = map[ErrTyp]string{
	WrongFormat:   "WrongFormat",
	EncodingError: "EncodingErr",
	EmptyArray:    "EmptyError",
	WrongCommand:  "WrongCommand",
}

var DefaultErrMsgMap = map[ErrTyp]string{
	WrongFormat:   "Request is in invalid format",
	EncodingError: "An Error has occured during Encoding",
	EmptyArray:    "Request is an empty array",
	WrongCommand:  "Command is not valid",
}

func (r *RespSimpleError) Error() string {
	return fmt.Sprintf("%s, %s", ErrPrefixMap[r.Prefix], r.Message)
}
