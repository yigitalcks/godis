package resp

import (
	"bufio"
	"strconv"
)

func (r *RespSimpleString) Encode(w *bufio.Writer) error {
	// +<data>\r\n

	// +
	err := w.WriteByte(PrefixSimpleString[0])
	if err != nil {
		return NewRespErr(EncodingError, err.Error())
	}
	// <data>
	_, err = w.WriteString(r.Value)
	if err != nil {
		return NewRespErr(EncodingError, err.Error())
	}
	// \r\n
	_, err = w.Write(CRLF)
	if err != nil {
		return NewRespErr(EncodingError, err.Error())
	}

	return nil
}

func (r *RespSimpleError) Encode(w *bufio.Writer) error {
	// -<data>\r\n

	// -
	err := w.WriteByte(PrefixSimpleError[0])
	if err != nil {
		return NewRespErr(EncodingError, err.Error())
	}
	// <data>
	_, err = w.WriteString(r.Value)
	if err != nil {
		return NewRespErr(EncodingError, err.Error())
	}
	// \r\n
	_, err = w.Write(CRLF)
	if err != nil {
		return NewRespErr(EncodingError, err.Error())
	}

	return nil
}

func (r *RespInteger) Encode(w *bufio.Writer) error {
	// :[<+|->]<value>\r\n

	// :
	err := w.WriteByte(PrefixInteger[0])
	if err != nil {
		return NewRespErr(EncodingError, err.Error())
	}
	// <value>
	_, err = w.WriteString(strconv.Itoa(r.Value))
	if err != nil {
		return NewRespErr(EncodingError, err.Error())
	}
	// \r\n
	_, err = w.Write(CRLF)
	if err != nil {
		return NewRespErr(EncodingError, err.Error())
	}

	return nil
}

func (r *RespBulkString) Encode(w *bufio.Writer) error {
	// $<length>\r\n<data>\r\n

	// Null Bulk String
	if len(r.Value) == 0 {
		// $
		if err := w.WriteByte(PrefixBulkString[0]); err != nil {
			return NewRespErr(EncodingError, err.Error())
		}
		// -1
		if _, err := w.WriteString(strconv.Itoa(-1)); err != nil {
			return NewRespErr(EncodingError, err.Error())
		}
		// \r\n
		if _, err := w.Write(CRLF); err != nil {
			return NewRespErr(EncodingError, err.Error())
		}
		return nil
	}

	// Non-Null Bulk String

	// $
	if err := w.WriteByte(PrefixBulkString[0]); err != nil {
		return NewRespErr(EncodingError, err.Error())
	}
	// <length>
	if _, err := w.WriteString(strconv.Itoa(len(r.Value))); err != nil {
		return NewRespErr(EncodingError, err.Error())
	}
	// \r\n
	if _, err := w.Write(CRLF); err != nil {
		return NewRespErr(EncodingError, err.Error())
	}
	// <data>
	if _, err := w.Write(r.Value); err != nil {
		return NewRespErr(EncodingError, err.Error())
	}
	// \r\n
	if _, err := w.Write(CRLF); err != nil {
		return NewRespErr(EncodingError, err.Error())
	}

	return nil
}

func (r *RespArray) Encode(w *bufio.Writer) error {
	// *<number-of-elements>\r\n<element-1>...<element-n>

	// *
	if err := w.WriteByte(PrefixArray[0]); err != nil {
		return NewRespErr(EncodingError, err.Error())
	}
	// <number-of-elements>
	if _, err := w.WriteString(strconv.Itoa(len(r.Value))); err != nil {
		return NewRespErr(EncodingError, err.Error())
	}
	// \r\n
	if _, err := w.Write(CRLF); err != nil {
		return NewRespErr(EncodingError, err.Error())
	}
	// <element-1>...<element-n>
	for _, element := range r.Value {
		if err := element.Encode(w); err != nil {
			return NewRespErr(EncodingError, err.Error())
		}
	}

	return nil
}
