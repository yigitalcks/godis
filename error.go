package main

import "fmt"

type ErrTyp int

const (
	WrongFormat ErrTyp = iota
)

var ErrPrefixMap = map[ErrTyp]string{
	WrongFormat: "WRONGFORMAT",
}

var DefaultErrMsgMap = map[ErrTyp]string{
	WrongFormat: "Request is in invalid format",
}

type RespErr struct {
	Prefix ErrTyp
	Msg    string
}

func (r *RespErr) Error() string {
	if len(r.Msg) == 0 {
		return fmt.Sprintf("%s, %s", ErrPrefixMap[r.Prefix], DefaultErrMsgMap[r.Prefix])
	}
	return fmt.Sprintf("%s, %s", ErrPrefixMap[r.Prefix], r.Msg)
}

func NewRespErr(prefix ErrTyp, msg string) *RespErr {
	return &RespErr{
		Prefix: prefix,
		Msg:    msg,
	}
}
