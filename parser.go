package main

import (
	"bytes"
	"strconv"
)

const MinBulkStringSize = 4 // in bytes

func ParseRequest(req []byte) ([][]byte, error) {

	if len(req) == 0 {
		return nil, NewRespErr(WrongFormat, "Too short request.")
	}

	i := bytes.Index(req, []byte("\r\n"))
	if i == -1 {
		return nil, NewRespErr(WrongFormat, "")
	}
	if len(req) < i+3 {
		return nil, NewRespErr(WrongFormat, "")
	}

	head := req[:i]
	if len(head) < 2 {
		return nil, NewRespErr(WrongFormat, "")
	}

	if head[0] != '*' {
		return nil, NewRespErr(WrongFormat, "")
	}

	num, err := strconv.ParseInt(string(head[1:]), 10, 64)
	if err != nil {
		return nil, NewRespErr(WrongFormat, "")
	}

	body := req[i+2:]
	comAndArgs, err := parseRequestData(int(num), body)
	if err != nil {
		return nil, err
	}

	return comAndArgs, nil
}

func parseRequestData(n int, s []byte) ([][]byte, error) {

	comAndArgs := make([][]byte, 0, n)
	for i := range n {
		if len(s) < MinBulkStringSize || s[0] != '$' {
			return nil, NewRespErr(WrongFormat, "")
		}

		idx := bytes.Index(s, []byte("\r\n"))
		if idx == -1 {
			return nil, NewRespErr(WrongFormat, "")
		}
		if len(s) < idx+3 {
			return nil, NewRespErr(WrongFormat, "")
		}

		head := s[:idx]
		if len(head) < 2 {
			return nil, NewRespErr(WrongFormat, "")
		}

		nBytes, err := strconv.ParseInt(string(s[1:idx]), 10, 64)
		if err != nil {
			return nil, NewRespErr(WrongFormat, "")
		}

		body := s[idx+2:]
		if len(body) < int(nBytes)+2 {
			return nil, NewRespErr(WrongFormat, "")
		}

		d := body[:nBytes]
		body = body[nBytes:]
		idxCLRF := bytes.Index(body, []byte("\r\n"))
		if idxCLRF == -1 {
			return nil, NewRespErr(WrongFormat, "")
		}

		if i < n-1 {
			if len(body[idxCLRF:]) < 3 {
				return nil, NewRespErr(WrongFormat, "")
			}
			s = body[idxCLRF+2:]
		}
		comAndArgs = append(comAndArgs, body[:idxCLRF])
	}
	return comAndArgs, nil
}
