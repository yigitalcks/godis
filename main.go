package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"time"

	"godis/internal/interpreter"
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

	fmt.Println("godis: starting on :6379")

	intr := interpreter.NewInterpreter(&data, expiryJobs)
	es := interpreter.NewExpiryScheduler(expiryJobs, &data)

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("godis: failed to bind to port 6379:", err)
		os.Exit(1)
	}
	defer l.Close()

	es.Start()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("godis: accept error:", err)
			os.Exit(1)
		}
		fmt.Println("godis: new connection from", conn.RemoteAddr())
		go handle_command(conn, intr)
	}
}

func writeResponse(c net.Conn, w *bufio.Writer, v resp.RespType) error {
	if err := c.SetWriteDeadline(time.Now().Add(5 * time.Second)); err != nil {
		return fmt.Errorf("set write deadline: %w", err)
	}
	if err := v.Encode(w); err != nil {
		return fmt.Errorf("encode: %w", err)
	}
	if err := w.Flush(); err != nil {
		if errors.Is(err, os.ErrDeadlineExceeded) {
			return fmt.Errorf("write deadline exceeded")
		}
		return fmt.Errorf("flush: %w", err)
	}
	return nil
}

func handle_command(c net.Conn, intr *interpreter.Interpreter) {
	defer c.Close()

	w := bufio.NewWriter(c)
	r := bufio.NewReaderSize(c, MaxRequestSize)

	buf := make([]byte, MaxRequestSize)
	for {
		err := c.SetReadDeadline(time.Now().Add(5 * time.Second))
		if err != nil {
			fmt.Println("Error setting read deadline: ", err.Error())
			return
		}

		n, err := r.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			} else if errors.Is(err, os.ErrDeadlineExceeded) {
				fmt.Println("godis: read deadline exceeded for", c.RemoteAddr())
			} else {
				fmt.Println("godis: read error from", c.RemoteAddr(), "-", err)
			}
			return
		}

		if r.Buffered() > 0 {
			fmt.Println("godis: request too large from", c.RemoteAddr())
			errResp := resp.NewRespSimpleError(resp.RequestTooLarge)
			if err = writeResponse(c, w, errResp); err != nil {
				fmt.Println("godis: write error to", c.RemoteAddr(), "-", err)
			}
			return
		}

		req := parser.ParseRequest(buf[:n])
		if _, isErr := req.(*resp.RespSimpleError); isErr {
			fmt.Println("godis: parse error from", c.RemoteAddr(), "-", req)
			if err = writeResponse(c, w, req); err != nil {
				fmt.Println("godis: write error to", c.RemoteAddr(), "-", err)
			}
			return
		}

		res := intr.Execute(req.(*resp.RespArray))
		if _, isErr := res.(*resp.RespSimpleError); isErr {
			fmt.Println("godis: command error from", c.RemoteAddr(), "-", res)
		}

		if err = writeResponse(c, w, res); err != nil {
			fmt.Println("godis: write error to", c.RemoteAddr(), "-", err)
			return
		}
	}
}

func handleExpiryEvents() {

}
