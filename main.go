package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"godis/internal/interpreter"
	"godis/internal/parser"
	"godis/internal/resp"
)

const BufSize = 2048

func main() {
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handle_command(conn)
	}
}

func handle_command(c net.Conn) {
	defer c.Close()

	w := bufio.NewWriterSize(c, BufSize)

	buf := make([]byte, BufSize)
	for {
		err := c.SetReadDeadline(time.Now().Add(5 * time.Second))
		if err != nil {
			fmt.Println("Error setting read deadline: ", err.Error())
			return
		}

		// TODO buffered reading
		n, err := c.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			} else if errors.Is(err, os.ErrDeadlineExceeded) {
				fmt.Println("Read deadline exceeded")
			} else {
				fmt.Println("Error reading from connection: ", err.Error())
			}
			return
		}

		req := parser.ParseRequest(buf[:n])
		if _, isErr := req.(*resp.RespSimpleError); isErr {
			fmt.Println("Error parsing request: ", req)
			return
		}
		res := interpreter.Execute(req.(*resp.RespArray))
		if _, isErr := res.(*resp.RespSimpleError); isErr {
			fmt.Println("Error executing command: ", res)
			return
		}

		err = c.SetWriteDeadline(time.Now().Add(5 * time.Second))
		if err != nil {
			fmt.Println("Error setting write deadline: ", err.Error())
			return
		}

		err = res.Encode(w)
		if err != nil {
			fmt.Println("Error encoding response: ", err.Error())
			return
		}

		err = w.Flush()
		if err != nil {
			if errors.Is(err, os.ErrDeadlineExceeded) {
				fmt.Println("Write deadline exceeded")
			} else {
				fmt.Println("Error flushing response: ", err.Error())
			}
			return
		}
	}
}
