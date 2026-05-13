package interpreter

import (
	"bytes"
	"godis/internal/resp"
	"strconv"
	"sync"
	"time"
)

var (
	OP_SET_NX   = []byte("NX")
	OP_SET_XX   = []byte("XX")
	OP_SET_IFEQ = []byte("IFEQ")
	OP_SET_IFNE = []byte("IFNE")

	OP_SET_GET = []byte("GET")

	OP_SET_EX = []byte("EX")
	OP_SET_PX = []byte("PX")
)

type Interpreter struct {
	data       *sync.Map
	expiryJobs chan Expiryjob
	commandMap map[string]func([]resp.RespBulkString) resp.RespType
}

func NewInterpreter(data *sync.Map, expiryJobs chan Expiryjob) *Interpreter {
	i := &Interpreter{
		data:       data,
		expiryJobs: expiryJobs,
	}
	i.commandMap = map[string]func([]resp.RespBulkString) resp.RespType{
		"PING": executePING,
		"ECHO": executeECHO,
		"SET":  i.executeSET,
		"GET":  i.executeGET,
	}
	return i
}

func (i *Interpreter) Execute(input *resp.RespArray) resp.RespType {

	if len(input.Value) == 0 {
		return resp.NewRespSimpleError(resp.EmptyArray)
	}

	cmd, ok := input.Value[0].(*resp.RespBulkString)
	if !ok {
		return resp.NewRespSimpleError(resp.WrongFormat)
	}

	handler, ok := i.commandMap[string(cmd.Value)]
	if !ok {
		return resp.NewRespSimpleError(resp.WrongCommand)
	}

	args := make([]resp.RespBulkString, 0, len(input.Value)-1)
	for _, v := range input.Value[1:] {
		bs, ok := v.(*resp.RespBulkString)
		if !ok {
			return resp.NewRespSimpleError(resp.WrongFormat)
		}
		// Null Bulk String e izin veriliyor.
		args = append(args, *bs)
	}

	return handler(args)
}

func executePING(args []resp.RespBulkString) resp.RespType {
	if len(args) > 1 {
		return resp.NewRespSimpleError(resp.WrongCommand)
	}

	if len(args) == 0 {
		return resp.NewRespSimpleString("PONG")
	} else {
		return resp.NewRespBulkString(args[0].Value)
	}
}

func executeECHO(args []resp.RespBulkString) resp.RespType {
	if len(args) != 1 {
		return resp.NewRespSimpleError(resp.WrongCommand)
	}
	return resp.NewRespBulkString(args[0].Value)
}

// Key: string
// Value: []byte
func (i *Interpreter) executeSET(args []resp.RespBulkString) resp.RespType {
	if len(args) < 2 {
		return resp.NewRespSimpleError(resp.WrongCommand)
	}

	key := args[0]
	if len(key.Value) == 0 {
		return resp.NewRespSimpleError(resp.WrongCommand, "Key cannot be null")
	}

	var expiry time.Duration

	opts := args[2:]
	for idx := 0; idx < len(opts); idx++ {
		op := opts[idx].Value

		switch {
		case bytes.Equal(op, OP_SET_EX), bytes.Equal(op, OP_SET_PX):
			unit := time.Second
			if bytes.Equal(op, OP_SET_PX) {
				unit = time.Millisecond
			}
			idx++
			if idx >= len(opts) {
				return resp.NewRespSimpleError(resp.WrongCommand)
			}
			n, err := strconv.Atoi(string(opts[idx].Value))
			if err != nil || n <= 0 {
				return resp.NewRespSimpleError(resp.WrongCommand)
			}
			expiry = time.Duration(n) * unit

		// TODO: NX, XX, GET, IFEQ, IFNE

		default:
			return resp.NewRespSimpleError(resp.WrongCommand)
		}
	}

	if expiry > 0 {
		select {
		case i.expiryJobs <- *NewExpiryJob(string(key.Value), expiry):
		default:
			return resp.NewRespSimpleError(resp.WrongCommand)
		}
	}

	i.data.Store(string(key.Value), args[1].Value)
	return resp.NewRespSimpleString("OK")
}

func (i *Interpreter) executeGET(args []resp.RespBulkString) resp.RespType {
	if len(args) != 1 {
		return resp.NewRespSimpleError(resp.WrongCommand)
	}

	val, ok := i.data.Load(string(args[0].Value))
	if !ok {
		return resp.NewRespBulkString(nil)
	}

	return resp.NewRespBulkString(val.([]byte))
}
