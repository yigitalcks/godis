package main

import (
	"bufio"
	"errors"
	"io"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"godis/internal/interpreter"
	"godis/internal/logging"
	"godis/internal/parser"
	"godis/internal/resp"
)

const (
	// MaxRequestSize is the maximum number of bytes allowed per request.
	MaxRequestSize          = 512 * 1024 // 512 KB
	MaxConcurrentExpiryJobs = 1024
)

func main() {

	var data sync.Map
	expiryJobs := make(chan interpreter.Expiryjob, MaxConcurrentExpiryJobs)

	logging.Println("godis: starting on :6379")

	intr := interpreter.NewInterpreter(&data, expiryJobs)
	es := interpreter.NewExpiryScheduler(expiryJobs, &data)

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		logging.Println("godis: failed to bind to port 6379:", err)
		os.Exit(1)
	}
	defer l.Close()

	es.Start()

	for {
		conn, err := l.Accept()
		if err != nil {
			logging.Println("godis: accept error:", err)
			os.Exit(1)
		}
		logging.Println("godis: new connection from", conn.RemoteAddr())
		go handle_command(conn, intr)
	}
}

func writeResponse(c net.Conn, w *bufio.Writer, v resp.RespType) {
	if err := c.SetWriteDeadline(time.Now().Add(5 * time.Second)); err != nil {
		logging.Printf("set write deadline: %w", err)
		return
	}
	if err := v.Encode(w); err != nil {
		logging.Printf("encode: %w", err)
		return
	}
	if err := w.Flush(); err != nil {
		if errors.Is(err, os.ErrDeadlineExceeded) {
			logging.Printf("write deadline exceeded")
			return
		}
		logging.Printf("flush: %w", err)
	}
}

func commandText(req *resp.RespArray) string {
	if len(req.Value) == 0 {
		return "<empty>"
	}

	parts := make([]string, 0, len(req.Value))
	for _, value := range req.Value {
		bulk, ok := value.(*resp.RespBulkString)
		if !ok {
			parts = append(parts, "<invalid>")
			continue
		}
		if bulk.Value == nil {
			parts = append(parts, "<nil>")
			continue
		}
		parts = append(parts, string(bulk.Value))
	}
	return strings.Join(parts, " ")
}

func handle_command(c net.Conn, intr *interpreter.Interpreter) {
	defer c.Close()

	w := bufio.NewWriter(c)
	r := bufio.NewReaderSize(c, MaxRequestSize)

	buf := make([]byte, MaxRequestSize)
	for {
		err := c.SetReadDeadline(time.Now().Add(5 * time.Second))
		if err != nil {
			logging.Println("Error setting read deadline:", err.Error())
			return
		}

		n, err := r.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			} else if errors.Is(err, os.ErrDeadlineExceeded) {
				logging.Println("godis: read deadline exceeded for", c.RemoteAddr())
			} else {
				logging.Println("godis: read error from", c.RemoteAddr(), "-", err)
			}
			return
		}

		if r.Buffered() > 0 {
			logging.Println("godis: request too large from", c.RemoteAddr())
			errResp := resp.NewRespSimpleError(resp.RequestTooLarge)
			writeResponse(c, w, errResp)
			return
		}

		req := parser.ParseRequest(buf[:n])
		if _, isErr := req.(*resp.RespSimpleError); isErr {
			logging.Println("godis: parse error from", c.RemoteAddr(), "-", req)
			writeResponse(c, w, req)
			return
		}

		reqArray := req.(*resp.RespArray)
		res := intr.Execute(reqArray)
		if _, isErr := res.(*resp.RespSimpleError); isErr {
			logging.Println("godis: command executed from", c.RemoteAddr(), "-", commandText(reqArray), "- error:", res)
		} else {
			logging.Println("godis: command executed from", c.RemoteAddr(), "-", commandText(reqArray))
		}

		writeResponse(c, w, res)
	}
}
