package resp

import (
	"bufio"
	"strconv"
)

// TODO Null Checks

func (r *RespSimpleString) Encode(w *bufio.Writer) error {
	w.WriteByte(PrefixSimpleString[0])
	w.WriteString(r.Value)
	w.Write(CRLF)

	return nil
}

func (r *RespSimpleError) Encode(w *bufio.Writer) error {
	w.WriteByte(PrefixSimpleError[0])
	w.WriteString(r.Value)
	w.Write(CRLF)

	return nil
}

func (r *RespInteger) Encode(w *bufio.Writer) error {
	w.WriteByte(PrefixInteger[0])
	w.WriteString(strconv.Itoa(r.Value))
	w.Write(CRLF)

	return nil
}

func (r *RespBulkString) Encode(w *bufio.Writer) error {
	//TODO To be Implemented
	return nil
}

func (r *RespArray) Encode(w *bufio.Writer) error {
	//TODO To be Implemented
	return nil
}
