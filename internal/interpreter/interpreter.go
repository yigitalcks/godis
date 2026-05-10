package interpreter

import "godis/internal/resp"

var CommandHandlerMap = map[string]func([]resp.RespBulkString) resp.RespType{
	"PING": executePING,
	"ECHO": executeECHO,
}

func Execute(input *resp.RespArray) resp.RespType {

	if len(input.Value) == 0 {
		return resp.NewRespSimpleError(resp.EmptyArray)
	}

	cmd, ok := input.Value[0].(*resp.RespBulkString)
	if !ok {
		return resp.NewRespSimpleError(resp.WrongFormat)
	}

	handler, ok := CommandHandlerMap[string(cmd.Value)]
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
	if len(args) == 0 {
		return resp.NewRespSimpleString("PONG")
	}
	return &resp.RespBulkString{Value: args[0].Value}
}

func executeECHO(args []resp.RespBulkString) resp.RespType {
	if len(args) != 1 {
		return resp.NewRespSimpleError(resp.WrongCommand)
	}
	return &resp.RespBulkString{Value: args[0].Value}
}
