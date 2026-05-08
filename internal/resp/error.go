package resp

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

func NewRespErr(prefix ErrTyp, msg ...string) *RespErr {
	var m string
	if len(msg) > 0 {
		m = msg[0]
	}
	return &RespErr{
		Prefix: prefix,
		Msg:    m,
	}
}
