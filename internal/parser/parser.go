package parser

import (
	"bytes"
	"godis/internal/resp"
	"strconv"
)

const MinArraySize = 4 // in bytes

var parsers = map[byte]func([]byte) ([]byte, []byte, error){
	resp.PrefixSimpleString[0]: parseSimpleString,
	resp.PrefixSimpleError[0]:  parseSimpleError,
	resp.PrefixInteger[0]:      parseInteger,
	resp.PrefixBulkString[0]:   parseBulkString,
}

func ParseArray(s []byte) ([][]byte, error) {

	var comAndArgs [][]byte

	if len(s) < MinArraySize {
		return nil, resp.NewRespErr(resp.WrongFormat)
	}

	idx := bytes.Index(s, resp.PrefixArray)
	if idx == -1 || idx != 0 {
		return nil, resp.NewRespErr(resp.WrongFormat)
	}

	idx = bytes.Index(s, resp.CRLF)
	if idx == -1 || idx == 1 {
		return nil, resp.NewRespErr(resp.WrongFormat)
	}

	nElements, err := strconv.Atoi(string(s[1:idx]))
	if err != nil {
		return nil, resp.NewRespErr(resp.WrongFormat)
	}

	if nElements == 0 {
		return [][]byte{}, nil
	}

	if idx+2 >= len(s) {
		return nil, resp.NewRespErr(resp.WrongFormat)
	}
	s = s[idx+2:]

	for range nElements {
		if len(s) == 0 {
			return nil, resp.NewRespErr(resp.WrongFormat)
		}

		parser, ok := parsers[s[0]]
		if !ok {
			return nil, resp.NewRespErr(resp.WrongFormat)
		}

		val, sNew, err := parser(s)
		if err != nil {
			return nil, err
		}

		s = sNew
		comAndArgs = append(comAndArgs, val)
	}

	return comAndArgs, nil
}

func parseSimpleString(s []byte) ([]byte, []byte, error) {
	// TODO To be Implemented
	return nil, nil, nil
}

func parseSimpleError(s []byte) ([]byte, []byte, error) {
	// TODO To be Implemented
	return nil, nil, nil
}

func parseInteger(s []byte) ([]byte, []byte, error) {
	// TODO To be Implemented
	return nil, nil, nil
}

func parseBulkString(s []byte) ([]byte, []byte, error) {

	idx := bytes.Index(s, resp.CRLF)
	if idx < 2 { // in case of -1, 0 and 1
		return nil, nil, resp.NewRespErr(resp.WrongFormat)
	}

	valLen, err := strconv.Atoi(string(s[1:idx]))
	if err != nil {
		return nil, nil, resp.NewRespErr(resp.WrongFormat)
	}
	if valLen == -1 {
		return nil, s[idx+2:], nil
	}

	if len(s)-idx-4 < valLen {
		return nil, nil, resp.NewRespErr(resp.WrongFormat)
	}

	res := s[idx+2 : idx+2+valLen]
	s = s[idx+2+valLen:]

	idx = bytes.Index(s, resp.CRLF)
	if idx != 0 {
		return nil, nil, resp.NewRespErr(resp.WrongFormat)
	}

	return res, s[2:], nil
}
