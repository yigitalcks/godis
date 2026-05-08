package main

type Resp interface {
	Encode() []byte
	Decode([]byte) error
}

type SimpleStringResp struct {
	Value string
}

type SimpleErrorResp struct {
	Value string
}

type IntegerResp struct {
	Value int64
}

type BulkStringResp struct {
	Value string
}

type ArrayResp struct {
	Value []Resp
}
