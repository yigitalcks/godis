package resp

var (
	PrefixArray        = []byte{'*'}
	PrefixSimpleString = []byte{'+'}
	PrefixSimpleError  = []byte{'-'}
	PrefixInteger      = []byte{':'}
	PrefixBulkString   = []byte{'$'}

	CRLF = []byte{'\r', '\n'}
)

type RespSimpleString struct {
	Value string
}

func NewRespSimpleString(value string) *RespSimpleString {
	return &RespSimpleString{
		Value: value,
	}
}

type RespSimpleError struct {
	Prefix  ErrTyp
	Message []byte
}

func NewRespSimpleError(prefix ErrTyp, message ...string) *RespSimpleError {
	var msg []byte
	if len(message) > 0 {
		msg = []byte(message[0])
	} else {
		msg = []byte(DefaultErrMsgMap[prefix])
	}

	return &RespSimpleError{
		Prefix:  prefix,
		Message: msg,
	}
}

type RespInteger struct {
	Value int
}

func NewRespInteger(value int) *RespInteger {
	return &RespInteger{
		Value: value,
	}
}

type RespBulkString struct {
	Value []byte
}

func NewRespBulkString(value []byte) *RespBulkString {
	return &RespBulkString{
		Value: value,
	}
}

type RespArray struct {
	Value []RespType
}

func NewRespArray(value []RespType) *RespArray {
	return &RespArray{
		Value: value,
	}
}
