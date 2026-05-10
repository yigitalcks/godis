package parser

import (
	"bytes"
	"godis/internal/resp"
	"strconv"
)

const MinArraySize = 4 // in bytes

// Request should be an array of bulk strings
func ParseRequest(s []byte) resp.RespType {

	var array resp.RespArray

	if len(s) < MinArraySize {
		return resp.NewRespSimpleError(resp.WrongFormat)
	}

	idx := bytes.Index(s, resp.PrefixArray)
	if idx == -1 || idx != 0 {
		return resp.NewRespSimpleError(resp.WrongFormat)
	}

	idx = bytes.Index(s, resp.CRLF)
	if idx == -1 || idx == 1 {
		return resp.NewRespSimpleError(resp.WrongFormat)
	}

	nElements, err := strconv.Atoi(string(s[1:idx]))
	if err != nil {
		return resp.NewRespSimpleError(resp.WrongFormat)
	}

	if nElements == 0 {
		if len(s) == MinArraySize {
			return &resp.RespArray{}
		}
		return resp.NewRespSimpleError(resp.WrongFormat)
	}

	if idx+2 >= len(s) {
		return resp.NewRespSimpleError(resp.WrongFormat)
	}
	s = s[idx+2:]

	for range nElements {
		if len(s) == 0 {
			return resp.NewRespSimpleError(resp.WrongFormat)
		}

		if s[0] != resp.PrefixBulkString[0] {
			return resp.NewRespSimpleError(resp.WrongFormat)
		}

		val, sNew := parseBulkString(s)
		if _, isErr := val.(*resp.RespSimpleError); isErr {
			return val
		}

		s = sNew
		array.Value = append(array.Value, val)
	}

	return &array
}

func parseBulkString(s []byte) (resp.RespType, []byte) {

	idx := bytes.Index(s, resp.CRLF)
	if idx < 2 { // in case of -1, 0 and 1
		return resp.NewRespSimpleError(resp.WrongFormat), nil
	}

	valLen, err := strconv.Atoi(string(s[1:idx]))
	if err != nil {
		return resp.NewRespSimpleError(resp.WrongFormat), nil
	}
	if valLen == -1 {
		return resp.NewRespBulkString(nil), s[idx+2:]
	}

	if len(s)-idx-4 < valLen {
		return resp.NewRespSimpleError(resp.WrongFormat), nil
	}

	res := resp.NewRespBulkString(s[idx+2 : idx+2+valLen])

	s = s[idx+2+valLen:]

	idx = bytes.Index(s, resp.CRLF)
	if idx != 0 {
		return resp.NewRespSimpleError(resp.WrongFormat), nil
	}

	return res, s[2:]
}
